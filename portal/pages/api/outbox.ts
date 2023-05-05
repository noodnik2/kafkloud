// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import {getServiceAddr, SERVICENAME_COURIER} from "@/routes/info";
import axios from 'axios';

export default async function handler(req: NextApiRequest, res: NextApiResponse<any>) {
    const deliveryPkg = req.body;
    console.log(`deliveryPkg(${deliveryPkg})`)
    const serviceAddr = getServiceAddr(SERVICENAME_COURIER);
    const message = JSON.stringify(deliveryPkg.message);
    const topic = deliveryPkg.topic;
    const key = deliveryPkg.key;
    const url = `${serviceAddr}/produce/${topic}/${key}`
    console.log(`POST url(${url}) message(${message})`);

    const axiosHeaders = {
        headers: {
            'Content-Type': 'application/json'
        }
    }

    return axios
        .post(url, message, axiosHeaders)
        .then((response) => {
            console.log(`POST OK; status(${response.statusText})`);
            res.status(200).json({status: response.statusText});
        })
        .catch((error) => {
            console.log(`POST received error(${error})`);
            res.status(500).json({error: error});
        });
}
