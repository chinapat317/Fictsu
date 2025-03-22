"use client"

import useSWR from "swr"
import { useState } from "react"
import { User } from "@/types/types"

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then(res => res.json())

export default function ChapterActions({ fiction_id, chapter_id, contributor_id }: {
    fiction_id: number,
    chapter_id: number,
    contributor_id: number,
}) {
    const [loading, setLoading] = useState(false)

    const { data } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/user`, fetcher)
    const user: User | null = data?.User_Profile || null

    if (!user || user.id !== contributor_id) {
        return null
    }

    const handleEdit = () => {
        window.location.href = `/f/${fiction_id}/${chapter_id}/edit`
    }

    const handleDelete = async () => {
        if (!confirm("Are you sure you want to delete this chapter?")) {
            return
        }

        setLoading(true)
        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/${chapter_id}/d`, {
                method: "DELETE",
                credentials: "include",
            })

            if (!res.ok) {
                throw new Error("Failed to delete chapter")
            }

            alert("Chapter deleted successfully")
            window.location.reload()
        } catch (error) {
            console.error("Error deleting chapter:", error)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="flex gap-2">
            <button onClick={handleEdit} className="text-blue-500 hover:underline" disabled={loading}>
                Edit
            </button>
            <button onClick={handleDelete} className="text-red-500 hover:underline" disabled={loading}>
                Delete
            </button>
        </div>
    )
}
