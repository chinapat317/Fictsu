"use client"

import useSWR from "swr"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { ChapterForm } from "@/types/types"
import { use, useState, useEffect } from "react"

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then((res) => res.json())

export default function EditChapterPage({ params }: { params: Promise<{ fiction_id: string; chapter_id: string }> }) {
    const router = useRouter()

    const [loading, setLoading] = useState(false)
    const { fiction_id, chapter_id } = use(params)
    const { register, handleSubmit, reset, formState: { errors } } = useForm<ChapterForm>()

    const { data: userData, error: userError } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/user`, fetcher)
    const { data: fictionData, error: fictionError } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}`, fetcher)
    const { data: chapterData, error: chapterError, mutate } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/${chapter_id}`, fetcher)

    useEffect(() => {
        if (userError || fictionError || chapterError) {
            alert("Failed to load data. Redirecting...")
            router.push("/")
            return
        }

        if (!userData || !fictionData?.Fiction || !chapterData) {
            return
        }

        if (!userData.User_Profile?.id) {
            alert("You must be logged in. Redirecting...")
            router.push("/")
            return
        }

        if (userData.User_Profile.id !== fictionData.Fiction.contributor_id) {
            alert("You are not the contributor. Redirecting...")
            router.push("/")
            return
        }

        reset({
            title: chapterData.title,
            content: chapterData.content,
        })
    }, [userData, fictionData, chapterData, reset, router, userError, fictionError, chapterError])

    const onSubmit = async (formData: ChapterForm) => {
        setLoading(true)
        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/${chapter_id}/u`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(formData),
        })

        setLoading(false)
        if (!response.ok) {
            alert("Failed to update chapter.")
            return;
        }

        alert("Chapter updated successfully!")
        mutate()
        router.push(`/f/${fiction_id}/${chapter_id}`)
    }

    if (!userData || !fictionData || !chapterData) {
        return <p className="text-center mt-10">Loading...</p>
    }

    return (
        <div className="max-w-2xl mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">Edit Chapter</h1>
            <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
                <input 
                    {...register("title", { required: true })} 
                    placeholder="Title *Required" 
                    className="w-full p-2 border rounded" 
                    defaultValue={chapterData?.title} 
                />
                {errors.title && <span className="text-red-500">Title is required</span>}

                <textarea 
                    {...register("content", { required: true })} 
                    placeholder="Content *Required" 
                    className="w-full p-2 border rounded" 
                    defaultValue={chapterData?.content} 
                />
                {errors.content && <span className="text-red-500">Content is required</span>}

                <button type="submit" className="bg-blue-500 text-white py-2 rounded-md" disabled={loading}>
                    {loading ? "Updating..." : "Save Changes"}
                </button>
            </form>
        </div>
    )
}
