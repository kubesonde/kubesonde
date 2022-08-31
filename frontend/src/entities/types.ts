export interface Dict {
    [key: string]: string
}
export interface IntDict {
    [key: string]: number
}

export interface DictDict {
    [key: string]: Dict
}
export interface BoolDict {
    [key: string]: boolean
}
