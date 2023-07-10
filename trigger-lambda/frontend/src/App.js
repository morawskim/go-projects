import './App.css';
import React from 'react';
import validator from '@rjsf/validator-ajv8';
import Form from '@rjsf/core';

const schema = {
    title: 'S3/PutObject',
    type: 'object',
    required: ['BUCKET_NAME', 'OBJECT_KEY'],
    properties: {
        BUCKET_NAME: { type: 'string', title: 'Bucket name', default: '' },
        OBJECT_KEY: { type: 'string', title: 'Object key', default: '' },
    },
};

const JsonSchemasForEvents = {
    's3/put-object': schema,
}

const TEMPLATES = {
    's3/put-object': `{"Records":[{"eventVersion":"2.0","eventSource":"aws:s3","awsRegion":"eu-central-1","eventTime":"1970-01-01T00:00:00.000Z","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"EXAMPLE"},"requestParameters":{"sourceIPAddress":"127.0.0.1"},"responseElements":{"x-amz-request-id":"EXAMPLE123456789","x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH"},"s3":{"s3SchemaVersion":"1.0","configurationId":"testConfigRule","bucket":{"name":"<<BUCKET_NAME>>","ownerIdentity":{"principalId":"EXAMPLE"},"arn":"arn:aws:s3:::data.linker.shop"},"object":{"key":"<<OBJECT_KEY>>","size":1024,"eTag":"0123456789abcdef0123456789abcdef","sequencer":"0A1B2C3D4E5F678901"}}}]}`,
}

function App() {
    const [formData, setFormData] = React.useState({});
    const [eventName, setEventName] = React.useState('');

    return (
        <div className="container">
            <div className="row">
                <div className="col">
                    <h1>Create AWS event payload</h1>

                    <div className="mb-3">
                        <label htmlFor="event" className="form-label">Email address</label>
                        <select defaultValue="" id="event" className="form-select" aria-label="Default select example"
                                onChange={event => {
                                    setEventName(event.target.value);
                                    setFormData({});
                                }}>
                            <option value="">Select event</option>
                            <optgroup label="S3">
                                <option value="s3/put-object">PutObject</option>
                            </optgroup>
                        </select>
                    </div>

                    <Form
                        className={'mb-3'}
                        schema={JsonSchemasForEvents[eventName] ?? {}}
                        validator={validator}
                        formData={formData}
                        onChange={(e) => setFormData(e.formData)}
                        uiSchema={{
                            'ui:submitButtonOptions': {
                                'norender': true,
                            },
                        }}
                    />

                    <button id="btnGenerate" type="button" className="btn btn-primary" onClick={event => {
                        event.preventDefault();
                        event.stopPropagation();

                        const eventName = document.getElementById('event').value;
                        if (!TEMPLATES.hasOwnProperty(eventName)) {
                            return;
                        }

                        let json = TEMPLATES[eventName];

                        for (let k in formData) {
                            if (!formData.hasOwnProperty(k)) {
                                continue;
                            }

                            json = json.replaceAll(`<<${k}>>`, formData[k]);
                        }
                        document.getElementById('event-payload').innerHTML = JSON.stringify(JSON.parse(json), null, 2);
                    }}>Generate event payload
                    </button>
                    <pre id="event-payload" className="mt-3"></pre>
                </div>
            </div>
        </div>
    );
}

export default App;
