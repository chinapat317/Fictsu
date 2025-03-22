"use client"

import useSWR from "swr"
import Image from "next/image"
import { useForm } from "react-hook-form"
import { useRouter } from "next/navigation"
import { FictionForm } from "@/types/types"
import { use, useState, useEffect } from "react"

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then((res) => res.json())

export default function FictionEditPage({ params }: { params: Promise<{ fiction_id: string }> }) {
    const router = useRouter()

    const { fiction_id } = use(params)
    const [loading, setLoading] = useState(false)
    const [cover, setCover] = useState("/default-cover.png")
    const { register, handleSubmit, reset, formState: { errors } } = useForm<FictionForm>()
    const { data: userData, error: userError } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/user`, fetcher)
    const { data: fictionData, error: fictionError, mutate } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}`, fetcher)

    useEffect(() => {
        if (userError || fictionError) {
            alert("Failed to load data. Redirecting...")
            router.push("/")
            return
        }

        if (!userData || !fictionData?.Fiction) {
            return
        }

        if (!userData.User_Profile || !userData.User_Profile.id) {
            alert("You must be logged in. Redirecting...")
            router.push("/")
            return
        }

        if (userData.User_Profile.id !== fictionData.Fiction.contributor_id) {
            alert("You are not the contributor. Redirecting...")
            router.push("/")
            return
        }

        reset(fictionData.Fiction)
        if (fictionData.Fiction.cover) {
            setCover(fictionData.Fiction.cover)
        }
    }, [userData, fictionData, fictionError, userError, reset, router])

    const onSubmit = async (formData: FictionForm) => {
        setLoading(true)
        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/u`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(formData),
        })

        setLoading(false)
        if (!response.ok) {
            alert("Failed to update fiction.")
            return
        }

        alert("Fiction updated successfully!")
        mutate()
        router.push(`/f/${fiction_id}`)
    }

    if (!fictionData || !userData) {
        return <p className="text-center mt-10">Loading...</p>
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            <h1 className="text-3xl font-bold mb-4">Edit Fiction</h1>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
                    {loading ? "Updating..." : "Save Changes"}
                </button>
            </form>
        </div>
    )
}
