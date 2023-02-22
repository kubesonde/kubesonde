export async function downloadJSON(filename: string, content: string) {
    const json = JSON.stringify(content)
    const blob = new Blob([json], { type: 'application/json' });
    const href = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = href;
    link.download = filename + ".json";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

export async function downloadImage(filename: string, blob: Blob) {
    const href = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = href;
    link.download = filename + ".png";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}
