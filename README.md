# receipt-processor

This repository houses a simple API for processing receipts and calculating points according to a pre-determined set of rules.

To setup this project and run it properly locally, here are some steps you'll need to follow. (Note: These steps assume your Go developer environment is already setup, and the language is correctly installed.)

### Initial Setup

1. Pull the repo from GitHub.
2. Ensure you have a tool installed to make API requests, such as Postman or cURL.
3. Configure Postman or cURL to run the API with the following URL: http://localhost:8080
4. Once the repo has been pulled and configured, run 'go run main.go' in your terminal from the receipt-processor directory.

### Testing

1. To test either of the available routes, you'll first need a single receipt that matches the following format (there can be any number of items in the "items" slice):

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```
```text

2. This receipt should be pasted into the request body of the Process Receipts (POST) route.
3. Run the API using 'go run main.go' or your IDE's launch system.
4. Hit the Process Receipts route and copy the "id" it returns.
5. To test the Calculate Receipt Points route, replace the {id} field in the request path with the ID you copied from step 4, then hit the route.
6. Validate that the "points" value returned by the Calculate Receipt Points route matches what you expect.

### Endpoints

The following API routes are currently supported:

Name: Process Receipts<br>
Request Type: POST<br>
Path: /receipts/process<br>
Description: This route accepts a receipt in JSON format, attaches a UUID to it, and returns said UUID.<br>

Name: Calculate Receipt Points<br>
Request Type: GET<br>
Path: /receipts/{id}/points<br>
Description: This route accepts an ID passed into the {id} portion of the route path. With that ID, the route searches for a matching receipt and calculates the appropriate number of points to grant based on the receipt's attributes.<br>

