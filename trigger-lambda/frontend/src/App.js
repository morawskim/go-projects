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

const formRef = React.createRef();

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
            <button onClick={event => {
                fetch(`${window.APP.BASE_PATH}/invoke`, {
                    method: "POST",
                    "body": props.payload,
                    headers: {
                        "Content-Type": "application/json",
                    },
                })
            }} className={"py-2 px-3 bg-sky-600 hover:bg-sky-800 text-white text-sm font-semibold rounded-md shadow focus:outline-none"}>
                Invoke
            </button>
            <pre
                className="mt-3"
                onClick={event => {
                    let range = new Range();
                    range.setStart(event.target, 0);
                    range.setEnd(event.target, 1);
                    document.getSelection().removeAllRanges();
                    document.getSelection().addRange(range);
                }}
            >{props.payload}</pre>
        </>
            ;
    }

    return <></>;
}


const Nav = (props) => {
    const dots = [];
    for (let i = 1; i <= props.totalSteps; i += 1) {
        const isActive = props.currentStep === i;
        dots.push((
            <div
                key={`step-${i}`}
                className={`step ${isActive ? 'step__active' : ''}`}
                onClick={() => props.goToStep(i)}
            >{props.titles[i-1] ?? ''}</div>
        ));
    }

    return (
        <div className={'step-container mb-3'}>{dots}</div>
    );
};

function App() {
    const [formData, setFormData] = React.useState({});
    const [eventName, setEventName] = React.useState('');
    const [payload, setPayload] = React.useState('');
    const [instance, setInstance] = React.useState(null);

    return (
        <div className="bg-slate-50 h-screen">
            <div className="container mx-auto">
                <div className="columns-1">
                    <h1 className={'text-center text-3xl py-2'}>Create AWS event payload</h1>

                    <StepWizard instance={setInstance} nav={<Nav titles={['Select event', 'Event parameters', 'Payload']} />}>
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
                                className={'mb-3 schema-form'}
                                schema={JsonSchemasForEvents[eventName] ?? {}}
                                validator={validator}
                                showErrorList={false}
                                formData={formData}
                                ref={formRef}
                                onChange={(e) => setFormData(e.formData)}
                                uiSchema={{
                                    'ui:submitButtonOptions': {
                                        'norender': true,
                                    },
                                }}
                            />
                            <button id="btnGenerate" type="button" className="float-right py-2 px-3 bg-sky-600 hover:bg-sky-800 text-white text-sm font-semibold rounded-md shadow focus:outline-none" onClick={event => {
                                event.preventDefault();
                                event.stopPropagation();
                                if (!formRef.current.validateForm()) {
                                    return;
                                }

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
