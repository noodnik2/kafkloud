import React, {useState, useEffect} from 'react';
import axios from 'axios';
import Head from "next/head";

export default function Consumer() {

    const [nConsumed, setNConsumed] = useState(-1);
    const [consumerUpdateCount, setConsumerUpdateCount] = useState(0);
    const [statusText, setStatusText] = useState('');

    const refreshConsumer = () => {
        // refresh consumer
        // const url = `${getServiceAddr("consumer")}`;
        const url = `/api/status`;
        console.log(`get(${url})`)
        axios.get(url)
            .then((response) => {
                const consumerStatusResponse = `consumerStatusResponse(${JSON.stringify(response.data)})`
                setNConsumed(response.data.n_consumed);
                setStatusText(consumerStatusResponse);
            })
            .catch((error) => {
                const axiosErrorStatus = `axiosErrorStatus(${error})`;
                console.error(axiosErrorStatus);
                setStatusText(axiosErrorStatus);
            });
    };

    useEffect(() => {
        // schedule a consumer refresh
        const timer = setTimeout(() => {
            // trigger a consumer refresh
            refreshConsumer();
        }, 2000);
        return () => clearTimeout(timer);
    }, [consumerUpdateCount]);

    const sendMessage = (message: FormDataEntryValue | null) => {
        console.log(`sendMessage(${message})`);
        const axiosHeaders = {
            headers: {
                'Content-Type': 'application/json'
            }
        }

        const url = `/api/outbox`;
        console.log(`post(${url})`)
        axios.post(url, message, axiosHeaders)
            .then((response) => {
                const producerRequestResponse = `producerRequestResponseStatusText(${response.statusText})`;
                setStatusText(producerRequestResponse);
                console.log(`triggering a consumer update`);
                setConsumerUpdateCount(consumerUpdateCount + 1);
            })
            .catch((error) => {
                const producerRequestError = `producerRequestError(${error})`;
                console.error(producerRequestError);
                setStatusText(producerRequestError);
            });
    };

    // @ts-ignore
    function submitProduceRequest(e) {
        e.preventDefault();
        const formData = new FormData(e.target);
        sendMessage(formData.get("message"));
    }

    const BUTTON_STYLES = "border border-blue-700 rounded m-1 py-2 px-4 bg-blue-500 hover:bg-blue-700 text-white font-bold"

    return (
        <main className="flex min-h-screen flex-col items-center justify-around p-24 font-sans backdrop ">
            <Head>
                <title>Consumer</title>
            </Head>

            <div className="border-2 bg-blue-50">
                <ConsumerCard nConsumed={nConsumed} />
            </div>

            <div className="border-2 bg-blue-50">
                {title("Producer Submission")}
                <div className="m-1 border-2">
                    <form method="post" onSubmit={submitProduceRequest}>
                        <div className="border-2">
                            <label>
                                <textarea
                                    className="mr-0.5"
                                    name="message"
                                    defaultValue='{"topic": "stream", "key":"k1", "message":{"text": "w1 w2"}}'
                                    rows={10}
                                    cols={75}
                                />
                            </label>
                        </div>
                        <div className="">
                            <label>
                                <button type="submit" className={BUTTON_STYLES}>Send</button>
                                <button type="reset" className={BUTTON_STYLES}>Reset</button>
                            </label>
                        </div>
                    </form>
                </div>
            </div>

            <div className="border-2">
                {title("Response")}
                <div className="p-1 m-0.5 bg-cyan-300 border-2 border-dotted">
                    <ResponseCard text={JSON.stringify(statusText)}></ResponseCard>
                </div>
            </div>

        </main>
    );
}

type ConsumerCardProps = {
    nConsumed: number
}

const ConsumerCard = ({nConsumed}: ConsumerCardProps) => {
    return (
        <>
            {title("Consumer Status")}
            <div className="m-1 border-2">
                <ul>
                    <li>nConsumed: {nConsumed}</li>
                </ul>
            </div>
        </>
    );
}

type ResponseCardProps = {
    text: string
}

const ResponseCard = ({text}: ResponseCardProps) => {
    return (
        <p>{text}</p>
    );
}

const title = (text: string) => {
    return (
        <h1 className="text-3xl underline decoration-2 font-bold bg-fuchsia-100">{text}</h1>
    );
}
