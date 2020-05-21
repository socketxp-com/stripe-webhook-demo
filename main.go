package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/stripe/stripe-go/webhook"
)

// application server port 
const port = ":8080"

func main() {
    secret := os.Getenv("STRIPE_SECRET")
    if secret == "" {
        fmt.Println("Required STRIPE_SECRET env variable is missing")
        os.Exit(1)
    }

    // incoming stripe webhook handler
    http.HandleFunc("/stripe", func(resp http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
            resp.WriteHeader(http.StatusBadRequest)
            return
        }

        // validating signature
        event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), secret)
        if err != nil {
            resp.WriteHeader(http.StatusBadRequest)
            fmt.Printf("Failed to validate signature: %s", err)
            return
        }

        switch event.Type {
        case "customer.subscription.created":
            // subscription create event
            customerID, ok := event.Data.Object["customer"].(string)
            if !ok {
                fmt.Println("customer key not found in event.Data.Object")
                return
            }

            subStatus, ok := event.Data.Object["status"].(string)
            if !ok {
                fmt.Println("status key not found in event.Data.Object")
                return
            }

            quantity, ok := event.Data.Object["quantity"]
            if !ok {
                fmt.Println("quantity key not found in event.Data.Object")
                return
            }

            fmt.Printf("customer %s subscription created, quantity (%f), current status: %s \n", customerID, quantity, subStatus)

        case "customer.subscription.updated":
            // subscription update event
            customerID, ok := event.Data.Object["customer"].(string)
            if !ok {
                fmt.Println("customer key not found in event.Data.Object")
                return
            }

            subStatus, ok := event.Data.Object["status"].(string)
            if !ok {
                fmt.Println("status key not found in event.Data.Object")
                return
            }

            quantity, ok := event.Data.Object["quantity"]
            if !ok {
                fmt.Println("quantity key not found in event.Data.Object")
                return
            }

            fmt.Printf("customer %s subscription updated, quantity(%f), current status: %s \n", customerID, quantity, subStatus)
        case "customer.subscription.deleted":
            // subscription deleted event 
            customerID, ok := event.Data.Object["customer"].(string)
            if !ok {
                fmt.Println("customer key not found in event.Data.Object")
                return
            }

            subStatus, ok := event.Data.Object["status"].(string)
            if !ok {
                fmt.Println("status key not found in event.Data.Object")
                return
            }

            quantity, ok := event.Data.Object["quantity"]
            if !ok {
                fmt.Println("quantity key not found in event.Data.Object")
                return
            }

            fmt.Printf("customer %s subscription deleted, quantity(%f), current status: %s \n", customerID, quantity, subStatus)
        default:
            fmt.Printf("Unknown event type received: %s\n", event.Type)
        }
    
    })

    fmt.Printf("Listening for Stripe webhooks on http://localhost%s/stripe \n", port)
    // starting the http server
    log.Fatal(http.ListenAndServe(port, nil))

}
