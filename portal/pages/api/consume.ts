// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import {getServiceAddr, SERVICENAME_MONITOR} from "@/routes/info";
import axios, {AxiosRequestConfig} from 'axios';

export default async function handler(req: NextApiRequest, outboundResponse: NextApiResponse<any>) {
    const serviceAddr = getServiceAddr(SERVICENAME_MONITOR);
    const url = `${serviceAddr}/consume?topics=${req.query['topics']}`

    // see: https://axios-http.com/docs/req_config
    const config: AxiosRequestConfig = {
        responseType: 'stream',
        headers: {
            'X-Requested-With': 'XMLHttpRequest',
            'Transfer-Encoding': 'chunked',
        },
    };
    await axios.get(url, config)
        .then((inboundResponse) => inboundResponse.data.pipe(outboundResponse))
        .catch((reason) => outboundResponse.status(424).send(`caught error(${reason})`))
}
