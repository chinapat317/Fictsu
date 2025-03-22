"use client"

import Image from "next/image"
import { useState } from "react"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { Fiction, FictionForm, ChapterForm } from "@/types/types"

export default function CreateFictionPage() {
    const router = useRouter()
    const chapterForm = useForm<ChapterForm>()

    const [loading, setLoading] = useState(false)
    const [cover, setCover] = useState("/default-cover.png")
    const [fiction, setFiction] = useState<Fiction | null>(null)
    const { register, handleSubmit, formState: { errors } } = useForm<FictionForm>()

    const onSubmitFiction = async (data: FictionForm) => {
        setLoading(true)
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/c`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(data),
                credentials: "include",
            })

            if (response.ok) {
                const fictionData = await response.json()
                setFiction(fictionData)
            } else {
                console.error("Failed to create fiction")
            }
        } catch (error) {
            console.error("Error submitting fiction form", error)
        } finally {
            setLoading(false)
        }
    }

    const onSubmitChapter = async (data: ChapterForm) => {
        if (!fiction) {
            return
        }
    
        setLoading(true)
        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction.id}/c`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(data),
                credentials: "include",
            })

            if (response.ok) {
                router.push(`/f/${fiction.id}`)
            } else {
                console.error("Failed to create chapter")
            }
        } catch (error) {
            console.error("Error submitting chapter form", error)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-4">Create Fiction</h1>
            <form onSubmit={handleSubmit(onSubmitFiction)} className="space-y-4">
                <div className="flex gap-6">
                    <Image
                        src={cover}
                        alt="Fiction Cover"
                        width={200}
                        height={300}
                        className="rounded-lg"
                    />
                    <div className="flex-1 space-y-2">
                        <input {...register("title", { required: true })} placeholder="Title *Required" className="w-full p-2 border rounded" />
                        {errors.title && <span className="text-red-500">Title is required</span>}

                        <input {...register("subtitle")} placeholder="Subtitle" className="w-full p-2 border rounded" />
                        <input {...register("author")} placeholder="Author" className="w-full p-2 border rounded" />
                        <input {...register("artist")} placeholder="Artist" className="w-full p-2 border rounded" />
                        <select {...register("status")} className="w-full p-2 border rounded">
                            <option value="Ongoing">Ongoing</option>
                            <option value="Completed">Completed</option>
                            <option value="Hiatus">Hiatus</option>
                            <option value="Dropped">Dropped</option>
                        </select>
                    </div>
                </div>
                <textarea {...register("synopsis", { required: true })} placeholder="Synopsis *Required" className="w-full p-2 border rounded" />
                {errors.synopsis && <span className="text-red-500">Synopsis is required</span>}

                <button type="submit" className="w-full bg-blue-500 text-white py-2 rounded" disabled={loading}>
                    {loading ? "Creating..." : "Create Fiction"}
                </button>
            </form>

            {fiction && (
                <div className="mt-8">
                    <h2 className="text-2xl font-semibold">Create First Chapter</h2>
                    <form onSubmit={chapterForm.handleSubmit(onSubmitChapter)} className="space-y-4 mt-4">
                        <input {...chapterForm.register("title", { required: true })} placeholder="Chapter Title *Required" className="w-full p-2 border rounded" />
                        {chapterForm.formState.errors.title && <span className="text-red-500">Title is required</span>}

                        <textarea {...chapterForm.register("content", { required: true })} placeholder="Chapter Content *Required" className="w-full p-2 border rounded" />
                        {chapterForm.formState.errors.content && <span className="text-red-500">Content is required</span>}

                        <button type="submit" className="w-full bg-green-500 text-white py-2 rounded" disabled={loading}>
                            {loading ? "Creating Chapter..." : "Create Chapter"}
                        </button>
                    </form>
                </div>
            )}
        </div>
    )
}
