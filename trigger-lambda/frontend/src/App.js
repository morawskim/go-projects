import './App.css';

const model = {};
const TEMPLATES = {
  "s3/put-object": `{"Records":[{"eventVersion":"2.0","eventSource":"aws:s3","awsRegion":"eu-central-1","eventTime":"1970-01-01T00:00:00.000Z","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"EXAMPLE"},"requestParameters":{"sourceIPAddress":"127.0.0.1"},"responseElements":{"x-amz-request-id":"EXAMPLE123456789","x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH"},"s3":{"s3SchemaVersion":"1.0","configurationId":"testConfigRule","bucket":{"name":"<<BUCKET_NAME>>","ownerIdentity":{"principalId":"EXAMPLE"},"arn":"arn:aws:s3:::data.linker.shop"},"object":{"key":"<<OBJECT_KEY>>","size":1024,"eTag":"0123456789abcdef0123456789abcdef","sequencer":"0A1B2C3D4E5F678901"}}}]}`,
}

function App() {
  return (
      <div className="container">
        <div className="row">
          <div className="col">
            <h1>Create AWS event payload</h1>

            <form>
              <div className="mb-3">
                <label htmlFor="event" className="form-label">Email address</label>
                <select defaultValue="" id="event" className="form-select" aria-label="Default select example" onChange={event => {
                  document.getElementById(`parameters-${event.target.value}`)?.classList?.toggle('d-none');
                }}>
                  <option value="">Select event</option>
                  <optgroup label="S3">
                    <option value="s3/put-object">PutObject</option>
                  </optgroup>
                </select>
              </div>

              <div id="parameters-s3/put-object" className="d-none payload-parameters" onChange={event => {
                if (event.target.tagName !== "INPUT") {
                  return;
                }

                model[event.target.name] = event.target.value;
              }}>
                <div className="mb-3">
                  <label htmlFor="txtBucket" className="form-label">Bucket</label>
                  <input type="text" className="form-control" id="txtBucket" name="BUCKET_NAME" />
                </div>

                <div className="mb-3">
                  <label htmlFor="txtObjectKey" className="form-label">Object key</label>
                  <input type="text" className="form-control" id="txtObjectKey" name="OBJECT_KEY" />
                </div>
              </div>

              <button id="btnGenerate" type="button" className="btn btn-primary" onClick={event => {
                event.preventDefault();
                event.stopPropagation();

                const eventName = document.getElementById('event').value;
                if (!TEMPLATES.hasOwnProperty(eventName)) {
                  return;
                }

                let json = TEMPLATES[eventName];

                for (let k in model) {
                  if (!model.hasOwnProperty(k)) {
                    continue;
                  }

                  json = json.replaceAll(`<<${k}>>`, model[k]);
                }
                document.getElementById('event-payload').innerHTML = JSON.stringify(JSON.parse(json), null, 2);
              }}>Generate event payload</button>
            </form>

            <pre id="event-payload" className="mt-3"></pre>
          </div>
        </div>
      </div>
  );
}

export default App;
