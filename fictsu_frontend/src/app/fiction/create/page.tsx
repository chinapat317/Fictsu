"use client"

import Image from "next/image"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"

var previewURL = "/images/fictsu_logo.png"

function ChangePreviewPic(e: React.ChangeEvent<HTMLInputElement>){
    var file = e.target.files?.[0] || null
    if(file) {
        var objectUrl = URL.createObjectURL(file)
        previewURL = objectUrl
    }
}


export default function CreateFictionPage() {
    const [cover, setCover] = useState<File | null>(null)
    const [title, setTitle] = useState<string>('')
    const [subtitle, setSubtitle] = useState<string>('')
    const [status, setStatus] = useState<string>('')
    const [synopsis, setSypnosis] = useState<string>('')
    const [author, setAuthor] = useState<string>('')
    const [artist, setArtist] = useState<string>('')
    return (
        <div className="max-w-4xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-4">Create Fiction</h1>
            <form onSubmit={() => {}} className="space-y-4">
                <div className="flex gap-6">
                <Image
                        src={previewURL}
                        alt="Fiction Cover"
                        width={200}
                        height={300}
                        className="rounded-lg"
                        unoptimized
                    />
                    <div className="flex-1 space-y-2">
                        <input type="text" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Title" className="w-full p-2 border rounded" required/>
                        <input type="text" value={subtitle} onChange={(e) => setSubtitle(e.target.value)} placeholder="Subtitle" className="w-full p-2 border rounded" required/>
                        <input type="text" value={author} onChange={(e) => setAuthor(e.target.value)} placeholder="Author" className="w-full p-2 border rounded" required/>
                        <input type="text" value={artist} onChange={(e) => setArtist(e.target.value)} placeholder="Artist" className="w-full p-2 border rounded" required/>
                        <select value={status} onChange={(e) => setStatus(e.target.value)}  className="w-full p-2 border rounded" required>
                            <option value="Ongoing">Ongoing</option>
                            <option value="Completed">Completed</option>
                            <option value="Hiatus">Hiatus</option>
                            <option value="Dropped">Dropped</option>
                        </select>
                        <button className="w-full bg-green-500 text-white py-2 rounded">
                        <input type="file"
                        accept="image/png, image/jpeg"
                        onChange={(e) => {
                            ChangePreviewPic(e)
                            setCover(e.target.files?.[0] || null)
                        }}
                        />
                        </button>
                    </div>
                </div>
                <textarea value={synopsis} placeholder="Synopsis" className="w-full p-2 border rounded" required/>
                <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded">
                </button>
            </form>
        </div>
    )
}