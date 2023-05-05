
export const SERVICENAME_MONITOR = "courier";
export const SERVICENAME_COURIER = "monitor";

export function getServiceAddr(serviceName: string): string {
    let serviceAddr;
    switch(serviceName) {
        case SERVICENAME_MONITOR: // TODO rename "consumer" to "monitor"
            serviceAddr = process.env.NEXT_PUBLIC_CONSUMER_ADDR;
            break;
        case SERVICENAME_COURIER:    // TODO rename "producer" to "courier"
            serviceAddr = process.env.NEXT_PUBLIC_PRODUCER_ADDR;
            break;
    }
    console.log(`for(${serviceName}) serviceAddr is(${serviceAddr})`)
    if (!serviceAddr) {
        throw `for(${serviceName}) serviceAddr is unknown`;
    }
    return serviceAddr;
};
