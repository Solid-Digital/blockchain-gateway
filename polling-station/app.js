const { join } = require("path");
const polka = require("polka");
var bodyParser = require("body-parser");
const { PORT = 3000 } = process.env;
const { HOST = "127.0.0.1" } = process.env;
const dir = join(__dirname, "public");
const serve = require("serve-static")(dir);
const request = require("request");

polka()
  .use(bodyParser.json(), serve)

  .get("/test/", (req, res) => {
    request.get("http://ec2-54-175-224-230.compute-1.amazonaws.com/").res();
  })

  .post("/vote/", (rx, tx) => {
    console.log(`received vote from ${rx.body.name} (${rx.body.fb_id})`);

    const options = {
      uri: "http://ec2-54-175-224-230.compute-1.amazonaws.com/adapter",
      method: "POST",
      json: {
        function: "vote",
        params: {
          option: parseInt(rx.body.option, 10),
          voterID: parseInt(rx.body.fb_id, 10)
        }
      },
      headers: {
        Accept: "application/json",
        Authorization: "Bearer 248dd60d-63d8-400e-abbc-f75906a909dd"
      }
    };

    request(options, function(err, res, body) {
      let respData = "";
      console.log(res.statusCode);

      if (res.statusCode == 200) {
        console.log(body);
        if (body == null) {
          respData = {
            result: "failure",
            message: `Transaction refused. User ${
              rx.body.fb_id
            } probably already voted previously`
          };
          tx.statusCode = 409;
        } else {
          respData = {
            result: "success",
            message: "Vote sucessfull!, Thank you for your vote."
          };
          tx.statusCode = 200;
        }
      } else {
        console.log(body);
        respData = {
          result: "failure",
          message: "The adapter returned an error. The vote failed to be send."
        };
        tx.statusCode = 500;
      }
      tx.end(JSON.stringify(respData));
    });
  })
  .listen(PORT, HOST, err => {
    if (err) throw err;
    console.log(`> Running on http://${HOST}:${PORT}`);
  });
