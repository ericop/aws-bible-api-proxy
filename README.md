# AWS Bible Api Proxy
Go on AWS Lambda for BibleScout.App

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://github.com/ericop/aws-bible-api-proxy)

> TODO: Split mp3 gathering to another repo for simplicity?

## How to get started locally and manually deploy
1. Generally follow https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html#golang-package-windows
2. Upload the `.zip` file on https://us-east-2.console.aws.amazon.com/lambda/home?region=us-east-2#/functions/bible-scout-proxy?tab=code
3. Test it at https://us-east-2.console.aws.amazon.com/lambda/home?region=us-east-2#/functions/bible-scout-proxy?tab=testing
  - **NOTE:** Sometimes it takes a few minutes for the new code to be the being tested, there is some kind of caching here
4. If testing works (or you want to test with the API Gateway semi integrated) go to https://us-east-2.console.aws.amazon.com/apigateway/home?region=us-east-2#/apis/dcu73qiiyi/resources/hal2db/methods/POST and chose the **Client "Test ⚡"** and enter the following in the body

   ```json
   {
       "urlText": "https%3A%2F%2Fdbt.io%2Faudio%2Fpath%3Fdam_id%3DENGESVN1DA%26book_id%3DActs%26v%3D2%26chapter_id%3D4",
       "code": "1nTh3B3g1nn1ngG0dCr34t3d"
   }
   ```

6. If this test worked it's time to deploy to the `default` stage at https://us-east-2.console.aws.amazon.com/apigateway/home?region=us-east-2#/apis/dcu73qiiyi/resources/hal2db and choose **Actions > Deploy API** > Deployment Stage: **default**
7. Now the ultimate test in **Postman**
- POST https://dcu73qiiyi.execute-api.us-east-2.amazonaws.com/default/bible-scout-proxy
- HEADERS with key:`x-api-key` value: `Genesis1-2InTheBeginningGodCreated` or `eAamcrnwum9yI7J9lDPYp3zLnDrBoqLcaLKBDDjc`
- BODY as raw JSON
   ```json
   {
       "urlText": "https%3A%2F%2Fdbt.io%2Faudio%2Fpath%3Fdam_id%3DENGESVN1DA%26book_id%3DActs%26v%3D2%26chapter_id%3D4",
       "code": "1nTh3B3g1nn1ngG0dCr34t3d"
   }
   ```
