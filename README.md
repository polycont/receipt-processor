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

`{__
  "retailer": "Target",__
  "purchaseDate": "2022-01-01",__
  "purchaseTime": "13:01",__
  "items": [__
    {__
      "shortDescription": "Mountain Dew 12PK",__
      "price": "6.49"__
    },{__
      "shortDescription": "Emils Cheese Pizza",__
      "price": "12.25"__
    },{__
      "shortDescription": "Knorr Creamy Chicken",__
      "price": "1.26"__
    },{__
      "shortDescription": "Doritos Nacho Cheese",__
      "price": "3.35"__
    },{__
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",__
      "price": "12.00"__
    }__
  ],__
  "total": "35.35"__
}`__

2. This receipt should be pasted into the request body of the Process Receipts (POST) route.
3. Run the API using 'go run main.go' or your IDE's launch system.
4. Hit the Process Receipts route and copy the "id" it returns.
5. To test the Calculate Receipt Points route, replace the {id} field in the request path with the ID you copied from step 4, then hit the route.
6. Validate that the "points" value returned by the Calculate Receipt Points route matches what you expect.

### Endpoints

The following API routes are currently supported:

Name: Process Receipts__
Request Type: POST__
Path: /receipts/process__
Description: This route accepts a receipt in JSON format, attaches a UUID to it, and returns said UUID.__

Name: Calculate Receipt Points__
Request Type: GET__
Path: /receipts/{id}/points__
Description: This route accepts an ID passed into the {id} portion of the route path. With that ID, the route searches for a matching receipt and calculates the appropriate number of points to grant based on the receipt's attributes.__

