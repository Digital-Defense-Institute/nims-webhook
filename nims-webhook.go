package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
    "os"

    "github.com/jomei/notionapi"
    "github.com/joho/godotenv"
)

var (
    client           *notionapi.Client
    assetsDatabaseID string 
    alertsDatabaseID string 
    authToken string
)

func init() {
    // load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

    // get env vars
    assetsDatabaseID = os.Getenv("NIMS_ASSETS_DATABASE_ID")
    alertsDatabaseID = os.Getenv("NIMS_ALERTS_DATABASE_ID")
    authToken = os.Getenv("NOTION_AUTH_TOKEN")

    // initialize notion client
    client = notionapi.NewClient(notionapi.Token(authToken))
}

func checkRelatedAssetExists(name string) (notionapi.ObjectID, error) {
    // search for the asset by title (name)
    filter := notionapi.PropertyFilter{
        Property: "Asset",
        RichText: &notionapi.TextFilterCondition{
            Equals: name,
        },
    }

    query := notionapi.DatabaseQueryRequest{
        Filter: &filter,
    }

    resp, err := client.Database.Query(context.Background(), notionapi.DatabaseID(assetsDatabaseID), &query)
    if err != nil {
        return "", err
    }

    // return the ID of the first matching item
    if len(resp.Results) > 0 {
        return notionapi.ObjectID(resp.Results[0].ID), nil
    }

    // no asset found, return empty PageID
    return "", nil
}

func createRelatedAsset(name, intIP string) (notionapi.ObjectID, error) {
    page := notionapi.PageCreateRequest{
        Parent: notionapi.Parent{DatabaseID: notionapi.DatabaseID(assetsDatabaseID)},
        Properties: notionapi.Properties{
            "Asset": notionapi.TitleProperty{
                Title: []notionapi.RichText{
                    {
                        Text: &notionapi.Text{Content: name},
                    },
                },
            },
            "Asset IP Address": notionapi.RichTextProperty{
                RichText: []notionapi.RichText{
                    {
                        Text: &notionapi.Text{Content: intIP},
                    },
                },
            },
        },
        Icon: &notionapi.Icon{
            Type: "external",
            External: &notionapi.FileObject{
                URL: "https://www.notion.so/icons/computer-chip_gray.svg",
            },
        },
    }

    resp, err := client.Page.Create(context.Background(), &page)
    if err != nil {
        return "", err
    }

    return resp.ID, nil
}

func formatTimestamp(timestamp string) (string, error) {
    // convert the timestamp to an integer
    ts, err := strconv.ParseInt(timestamp, 10, 64)
    if err != nil {
        return "", err
    }

    // convert milliseconds to seconds
    seconds := ts / 1000
    nanos := (ts % 1000) * 1000000

    // parse the timestamp
    t := time.Unix(seconds, nanos)

    // return the formatted timestamp
    return t.Format(time.RFC3339), nil
}

func addAlert(name, timestamp, hostname, intIP, link, details, metadata string) error {

    formattedTimestamp, err := formatTimestamp(timestamp)
    if err != nil {
        return fmt.Errorf("invalid timestamp format: %v", err)
    }

    // parse details and metadata
    var detailsJSON map[string]interface{}
    if err := json.Unmarshal([]byte(details), &detailsJSON); err != nil {
        return err
    }
    detailsFormatted, _ := json.MarshalIndent(detailsJSON, "", "    ")

    var metadataJSON map[string]interface{}
    if err := json.Unmarshal([]byte(metadata), &metadataJSON); err != nil {
        return err
    }
    metadataFormatted, _ := json.MarshalIndent(metadataJSON, "", "    ")

    // check or create asset
    assetID, err := checkRelatedAssetExists(hostname)
    if err != nil {
        return err
    }
    if assetID == "" {
        assetID, err = createRelatedAsset(hostname, intIP)
        if err != nil {
            return err
        }
    }

    // parse the formatted timestamp
    timeObj, err := time.Parse(time.RFC3339, formattedTimestamp)
    if err != nil {
        return err
    }
    notionDate := notionapi.Date(timeObj)
    notionDatePointer := &notionDate
    notionDateObject := &notionapi.DateObject{
        Start: notionDatePointer,
    }

    // create alert
    page := notionapi.PageCreateRequest{
        Parent: notionapi.Parent{DatabaseID: notionapi.DatabaseID(alertsDatabaseID)},
        Properties: notionapi.Properties{
            "Name": notionapi.TitleProperty{
                Title: []notionapi.RichText{
                    {
                        Text: &notionapi.Text{Content: name},
                    },
                },
            },
            "Alert Generated": notionapi.DateProperty{
                Date: notionDateObject,
            },
            "Details": notionapi.RichTextProperty{
                RichText: []notionapi.RichText{
                    {
                        Text: &notionapi.Text{Content: string(detailsFormatted)},
                    },
                },
            },
            "Metadata": notionapi.RichTextProperty{
                RichText: []notionapi.RichText{
                    {
                        Text: &notionapi.Text{Content: string(metadataFormatted)},
                    },
                },
            },
            "Affected Assets": notionapi.RelationProperty{
                Relation: []notionapi.Relation{
                    {ID: notionapi.PageID(assetID)},
                },
            },
            "Related URL": notionapi.URLProperty{
                URL: link,
            },
        },
        Icon: &notionapi.Icon{
            Type: "external",
            External: &notionapi.FileObject{
                URL: "https://www.notion.so/icons/bell_gray.svg",
            },
        },
    }

    _, err = client.Page.Create(context.Background(), &page)
    return err
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    // parse the incoming json
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, fmt.Sprintf("Error parsing webhook: %v", err), http.StatusBadRequest)
        return
    }

    // get name
    name := ""
    if val, ok := data["cat"].(string); ok {
        name = val
    }

    // get timestamp
    timestamp := ""
    var err error
    if val, ok := data["routing"].(map[string]interface{})["event_time"].(float64); ok {
        // Convert the float64 value to int64 to remove decimals
        timestamp = fmt.Sprintf("%d", int64(val))
    }

    // get hostname
    hostname := ""
    if val, ok := data["routing"].(map[string]interface{})["hostname"].(string); ok {
        hostname = val
    }

    // get ip address
    intIP := ""
    if val, ok := data["routing"].(map[string]interface{})["int_ip"].(string); ok {
        intIP = val
    }

    // get url
    link := ""
    if val, ok := data["link"].(string); ok {
        link = val
    }

    // turn details into json string
    detailsBytes, err := json.Marshal(data["detect"])
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshalling 'detect': %v", err), http.StatusInternalServerError)
        return
    }
    details := string(detailsBytes)

    // turn metadata into json string
    metadataBytes, err := json.Marshal(data["detect_mtd"])
    if err != nil {
        http.Error(w, fmt.Sprintf("Error marshalling 'detect_mtd': %v", err), http.StatusInternalServerError)
        return
    }
    metadata := string(metadataBytes)

    // add the alert
    err = addAlert(name, timestamp, hostname, intIP, link, details, metadata)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to add alert: %v", err), http.StatusInternalServerError)
        return
    }

    // success
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Alert added successfully")
}

func main() {
    // listen on port 9000 for webhook POST requests
    http.HandleFunc("/hooks/alert", webhookHandler)
    log.Println("Listening for webhooks on port 9000")
    log.Fatal(http.ListenAndServe(":9000", nil))
}
