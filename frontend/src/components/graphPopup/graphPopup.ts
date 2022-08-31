import tippy from "tippy.js"
import 'tippy.js/dist/tippy.css';

export const BuildPopup = (clientRect: DOMRect, message: string) => {
    let dummyDomEle = document.createElement('div');
    // @ts-ignore
    const popup = new tippy(dummyDomEle,{
        getReferenceClientRect: clientRect,
        trigger: 'manual',
        content: () => {
            let content = document.createElement('div');

            content.innerHTML = message;

            return content;
        }
    })
    return popup
}
/*let tip = new tippy(dummyDomEle, { // tippy props:
    getReferenceClientRect: ref.getBoundingClientRect, // https://atomiks.github.io/tippyjs/v6/all-props/#getreferenceclientrect
    trigger: 'manual', // mandatory, we cause the tippy to show programmatically.

    // your own custom props
    // content prop can be used when the target is a single element https://atomiks.github.io/tippyjs/v6/constructor/#prop
    content: () => {
        let content = document.createElement('div');

        content.innerHTML = 'Tippy content';

        return content;
    }
});

tip.show();*/
