import {linux} from './probesWrapped'

export interface NetstatInterface {
    protocol: string,
    local: {
        port: number,
        address: string
    },
    remote: {
        port: number,
        address: string
    }
    state: 'LISTEN' | 'ESTABLISHED' | 'TIME_WAIT' | 'SYN_SENT'
}
export function parseNetstat(rawOutput: string): NetstatInterface{
    const parser = linux()
    const result = parser(rawOutput,function(i: any){return i})
    return result
}
