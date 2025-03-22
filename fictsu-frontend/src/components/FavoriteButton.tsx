"use client"

import { useState, useEffect } from "react"

export default function FavoriteButton({ fiction_id }: { fiction_id: number }) {
    const [loading, setLoading] = useState(false)
    const [isFavorited, setIsFavorited] = useState(false)

    useEffect(() => {
        async function checkFavoriteStatus() {
            try {
                const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/fav/status`, {
                    credentials: "include",
                })

                if (res.status === 401) {
                    console.warn("User not logged in, skipping favorite status check.")
                    return
                }

                if (!res.ok) {
                    throw new Error("Failed to check favorite status")
                }

                const data = await res.json()
                setIsFavorited(data.is_favorited)
            } catch (error) {
                console.error("Error checking favorite status:", error)
            }
        }

        checkFavoriteStatus()
    }, [fiction_id])

    async function toggleFavorite() {
        setLoading(true)
        try {
            const method = isFavorited ? "DELETE" : "POST"
            const endpoint = isFavorited
            ? `${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/fav/rmv`
            : `${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/fav`

            const res = await fetch(endpoint, {
                method,
                credentials: "include",
            })

            if (res.status === 401) {
                alert("You need to log in to favorite a fiction.")
                setLoading(false)
                return
            }
    
            if (!res.ok) {
                throw new Error(`Failed to ${isFavorited ? "remove" : "add"} favorite`)
            }
            
            const data = await res.json()
            setIsFavorited(data.is_favorited)
        } catch (error) {
            console.error("Error toggling favorite:", error)
        } finally {
            setLoading(false)
        }
    }

    return (
        <button
            disabled={loading}
            onClick={toggleFavorite}
            className="p-2 text-red-500 hover:text-red-700"
            aria-label={isFavorited ? "Remove from favorites" : "Add to favorites"}
        >
            {isFavorited ? "❤️" : "🤍"}
        </button>
    )
}
