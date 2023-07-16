
export const SERVICENAME_MONITOR = "monitor";
export const SERVICENAME_COURIER = "courier";

export function getServiceAddr(serviceName: string): string {

    // console.log(`process.env:`)
    // Object.keys(process.env).forEach((prop, index, value)=> console.log(`${prop}=${value}`));

    let serviceAddr;
    switch(serviceName) {
        case SERVICENAME_MONITOR:
            serviceAddr = process.env.NEXT_PUBLIC_CONSUMER_ADDR;
            break;
        case SERVICENAME_COURIER:
            serviceAddr = process.env.NEXT_PUBLIC_PRODUCER_ADDR;
            break;
    }
    console.log(`for(${serviceName}) serviceAddr is(${serviceAddr})`)
    if (!serviceAddr) {
        throw `for(${serviceName}) serviceAddr is unknown`;
    }
    return serviceAddr;
};
