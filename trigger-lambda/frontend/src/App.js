import './App.css';
import React from 'react';
import validator from '@rjsf/validator-ajv8';
import Form from '@rjsf/core';
import Select from 'react-select'
import StepWizard from "react-step-wizard";

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
    's3/put-object': `{"Records":[{"eventVersion":"2.0","eventSource":"aws:s3","awsRegion":"eu-central-1","eventTime":"1970-01-01T00:00:00.000Z","eventName":"ObjectCreated:Put","userIdentity":{"principalId":"EXAMPLE"},"requestParameters":{"sourceIPAddress":"127.0.0.1"},"responseElements":{"x-amz-request-id":"EXAMPLE123456789","x-amz-id-2":"EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH"},"s3":{"s3SchemaVersion":"1.0","configurationId":"testConfigRule","bucket":{"name":"<<BUCKET_NAME>>","ownerIdentity":{"principalId":"EXAMPLE"},"arn":"arn:aws:s3:::<<BUCKET_NAME>>"},"object":{"key":"<<OBJECT_KEY>>","size":1024,"eTag":"0123456789abcdef0123456789abcdef","sequencer":"0A1B2C3D4E5F678901"}}}]}`,
}

const OPTIONS = [
    {value: '', label: 'Select event'},
    {
        label: 'S3',
        options: [
            {value: 's3/put-object', label: 'S3/PutObject'}
        ],
    },
]

function AwsPayloadEvent(props) {
    if (props.payload) {
        return <>
            <pre id="event-payload" className="mt-3">{props.payload}</pre>
            <button onClick={event => {
                fetch("/invoke", {
                    method: "POST",
                    "body": props.payload,
                    headers: {
                        "Content-Type": "application/json",
                    },
                })
            }} className={"btn btn-primary"}>Invoke</button>
        </>
            ;
    }

    return <></>;
}

function App() {
    const [formData, setFormData] = React.useState({});
    const [eventName, setEventName] = React.useState('');
    const [payload, setPayload] = React.useState('');
    const [instance, setInstance] = React.useState(null);

    return (
        <div className="container">
            <div className="row">
                <div className="col">
                    <h1>Create AWS event payload</h1>

                    <StepWizard instance={setInstance}>
                        <>
                            <div className="mb-3">
                                <label htmlFor="event" className="form-label">Email address</label>
                                <Select
                                    defaultValue={OPTIONS[0]}
                                    id={"event"}
                                    onChange={(event) => {
                                        setEventName(event.value);
                                        setFormData({});
                                        setPayload('');
                                        instance.nextStep();
                                    }}
                                    options={OPTIONS} />
                            </div>
                        </>

                        <>
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
                                setPayload(JSON.stringify(JSON.parse(json), null, 2))
                                instance.nextStep();
                            }}>Generate event payload
                            </button>
                        </>

                        <>
                            <AwsPayloadEvent payload={payload}></AwsPayloadEvent>
                        </>
                    </StepWizard>
                </div>
            </div>
        </div>
    );
}

export default App;
