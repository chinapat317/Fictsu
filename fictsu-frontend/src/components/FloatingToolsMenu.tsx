import { useState, useEffect } from "react"
import { Menu, X, ImagePlus, Bot } from "lucide-react"

export default function FloatingToolsMenu() {
    const [open, setOpen] = useState(false)
    const [loading, setLoading] = useState(false)
    const [AIResponse, setAIResponse] = useState("")
    const [inputPrompt, setInputPrompt] = useState("")

    const toggle_menu = () => setOpen(!open)

    const fetch_AI_response = async () => {
        if (!inputPrompt) {
            return
        }

        setLoading(true)
        setAIResponse("")

        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/ai/storyline/c`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ Message: inputPrompt }),
            })

            const data = await response.json()
            if (response.ok) {
                setAIResponse(data.Received_Message)
            } else {
                setAIResponse("Error: " + data.Error)
            }
        } catch (error) {
            setAIResponse("Failed to fetch AI response.")
            console.error("Error toggling favorite: ", error)
        }

        setLoading(false)
    }

    // Use useEffect to manipulate the DOM after the component has rendered
    useEffect(() => {
        const AI_chat_element = document.getElementById("AI-chat")
        if (AI_chat_element) {
            AI_chat_element.style.display = open ? "block" : "none"
        }
    }, [open]) // Effect runs when the "open" state changes

    return (
        <div className="fixed bottom-4 right-4">
            <div className="relative">
                <button
                    onClick={toggle_menu}
                    className="p-3 bg-gray-800 text-white rounded-full shadow-lg"
                >
                    {open ? <X size={24} /> : <Menu size={24} />}
                </button>

                {open && (
                    <div className="absolute bottom-12 right-0 w-48 bg-white shadow-lg rounded-lg p-2">
                        <h2 className="text-lg font-semibold">Tools</h2>
                        <button
                            className="flex items-center w-full p-2 hover:bg-gray-100 rounded"
                            onClick={() => document.getElementById("AI-chat")?.classList.toggle("hidden")}
                        >
                            <Bot size={20} className="mr-2" /> AI Chat
                        </button>
                        <button
                            className="flex items-center w-full p-2 hover:bg-gray-100 rounded"
                            onClick={() => console.log("Upload Image")}
                        >
                            <ImagePlus size={20} className="mr-2" /> Upload Image
                        </button>
                        <button
                            className="flex items-center w-full p-2 hover:bg-gray-100 rounded"
                            onClick={() => console.log("Generate Image")}
                        >
                            <ImagePlus size={20} className="mr-2" /> Generate Image
                        </button>
                    </div>
                )}
            </div>
            {open && (
                <div
                    className="fixed inset-0 bg-transparent"
                    onClick={() => setOpen(false)}
                ></div>
            )}

            <div id="AI-chat" className="hidden fixed bottom-20 right-4 bg-white shadow-lg p-4 rounded-lg w-80">
                <h2 className="text-lg font-semibold">AI Assistant</h2>
                <textarea
                    className="w-full p-2 border rounded mt-2"
                    placeholder="Describe what you need..."
                    value={inputPrompt}
                    onChange={(e) => setInputPrompt(e.target.value)}
                ></textarea>
                <button
                    className="w-full bg-blue-500 text-white p-2 rounded mt-2"
                    onClick={fetch_AI_response}
                    disabled={loading}
                >
                    {loading ? "Generating..." : "Generate"}
                </button>
                {AIResponse && <p className="mt-2 p-2 border rounded bg-gray-100">{AIResponse}</p>}
            </div>
        </div>
    )
}
