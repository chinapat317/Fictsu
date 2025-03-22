"use client"

import useSWR from "swr"
import Link from "next/link"
import Image from "next/image"
import { useRouter } from "next/navigation"
import { useState, useEffect } from "react"
import { User, Fiction } from "@/types/types"

const fetcher = (url: string) => fetch(url, { credentials: "include" }).then((res) => res.json())

export default function UserProfilePage() {
    const router = useRouter()

    const [activeTab, setActiveTab] = useState<"favorites" | "contributions">("favorites")
    const { data, error, mutate } = useSWR(`${process.env.NEXT_PUBLIC_BACKEND_API}/user`, fetcher)

    useEffect(() => {
        if (error || (data && !data.User_Profile)) {
            alert("Please login to access your profile.")
            router.push("/")
        }
    }, [error, data, router])

    if (!data || !data.User_Profile) {
        return <p className="text-center mt-10">Redirecting...</p>
    }

    const user: User = data.User_Profile
    return (
        <div className="max-w-3xl mx-auto p-6 bg-white shadow-md rounded-lg mt-10">
            <div className="flex items-center space-x-4">
                <img src={user.avatar_url} alt="Avatar" className="w-16 h-16 rounded-full" referrerPolicy="no-referrer" />
                <div>
                    <h1 className="text-2xl font-bold">{user.name}</h1>
                    <p className="text-gray-600">{user.email}</p>
                    <p className="text-sm text-gray-500">Joined: {new Date(user.joined).toLocaleDateString()}</p>
                </div>
            </div>
            <div className="mt-6 border-b flex">
                <button
                    onClick={() => setActiveTab("favorites")}
                    className={`px-4 py-2 text-lg font-medium ${
                        activeTab === "favorites"
                        ? "border-b-2 border-blue-600 text-blue-600"
                        : "text-gray-500"
                    }`}
                >
                    Favorites
                </button>
                <button
                    onClick={() => setActiveTab("contributions")}
                    className={`px-4 py-2 text-lg font-medium ${
                        activeTab === "contributions"
                        ? "border-b-2 border-blue-600 text-blue-600"
                        : "text-gray-500"
                    }`}
                >
                    Contributions
                </button>
            </div>
            <div className="mt-4">
                {activeTab === "favorites" ? (
                    <FictionList fictions={user.fav_fictions} emptyMessage="No favorite fictions yet." />
                ) : (
                    <FictionList
                        fictions={user.contributed_fic}
                        emptyMessage="No contributed fictions yet."
                        allowEdit={true}
                        mutate={mutate}
                    />
                )}
            </div>
        </div>
    )
}

function FictionList({ fictions, emptyMessage, allowEdit, mutate }: {
    fictions: Fiction[] | null | undefined,
    emptyMessage: string,
    allowEdit?: boolean,
    mutate?: () => void,
}) {
    const [fictionList, setFictionList] = useState<Fiction[]>(fictions || [])

    useEffect(() => {
        setFictionList(fictions || [])
    }, [fictions])

    if (!fictions || fictions.length === 0) {
        return <p className="text-gray-500 text-center">{emptyMessage}</p>
    }

    const handleEdit = (fiction_id: number) => {
        window.location.href = `/f/${fiction_id}/edit`
    }

    const handleDelete = async (fiction_id: number) => {
        if (!confirm("Are you sure you want to delete this fiction?")) {
            return
        }

        try {
            const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/d`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            })

            if (!response.ok) {
                throw new Error("Failed to delete fiction")
            }

            alert("Fiction deleted successfully")
            setFictionList((prev) => prev.filter((fic) => fic.id !== fiction_id))
            if (mutate) {
                await mutate()
            }
        } catch (error) {
            console.error(error)
        }
    }

    return (
        <ul className="space-y-4">
            {fictions.map((fiction) => (
                <li key={fiction.id} className="p-4 bg-gray-100 rounded-lg shadow-sm flex">
                    <Image 
                        src={fiction.cover || "/default-cover.png"} 
                        alt={fiction.title} 
                        width={60} 
                        height={90} 
                        className="rounded-md"
                    />
                    <div className="ml-4">
                        <Link href={`/f/${fiction.id}`}>
                            <h3 className="text-lg font-semibold hover:text-blue-600">{fiction.title}</h3>
                        </Link>
                        <p className="text-sm text-gray-600">{fiction.chapters?.length ?? 0} Chapters</p>
                        <p className="text-xs text-gray-500 mt-1">{new Date(fiction.created).toLocaleDateString()}</p>

                        {allowEdit && (
                            <div className="mt-2 space-x-2">
                                <button
                                    onClick={() => handleEdit(fiction.id)}
                                    className="px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
                                >
                                    Edit
                                </button>
                                <button
                                    onClick={() => handleDelete(fiction.id)}
                                    className="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
                                >
                                    Delete
                                </button>
                            </div>
                        )}
                    </div>
                </li>
            ))}
        </ul>
    )
}
