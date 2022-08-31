import {fireEvent, render, screen} from '@testing-library/react'
import '@testing-library/jest-dom'
import {ExampleGraphComponent} from "./ExampleGraph";

test('It Renders correct items when instantiated', async () => {
        render(<ExampleGraphComponent />)
        const graphTable = screen.getByRole('graphTable')
        expect(graphTable.children.length).toEqual(2) // Thead - Tbody
        const tbody: Element|null = graphTable.children.item(1)
        expect(tbody).not.toBeNull()
        expect(tbody?.children.item(0)).toHaveTextContent("Deployment-2")
        expect(tbody?.children.item(1)).toHaveTextContent("Deployment-1")
        const deployment2 = await screen.findByText("Deployment-2")
        const checkBox = deployment2.children[0] as HTMLInputElement
        expect(checkBox.checked).toBe(true)
        })
test('It handles deployment click', async () => {
    render(<ExampleGraphComponent />)
    const deployment2 = await screen.findByText("Deployment-2")
    const checkBox = deployment2.children[0] as HTMLInputElement
    fireEvent.click(checkBox)
    expect(checkBox.checked).toBe(false)
})

test('It handles enable click', async () => {
    render(<ExampleGraphComponent />)
    const checkboxes = await screen.findAllByRole("checkbox")
    expect(checkboxes.length).toEqual(8)
    const firstCheckbox = checkboxes[0] as HTMLInputElement
    expect(firstCheckbox.checked).toBe(true)
    fireEvent.click(firstCheckbox)
    expect(firstCheckbox.checked).toBe(false)
})
