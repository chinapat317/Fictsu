"use client"

import Link from "next/link"
import { User } from "@/types/types"
import { useEffect, useState, useRef } from "react"
import { useRouter, usePathname } from "next/navigation"

export default function Navbar() {
    const router = useRouter()
    const pathname = usePathname()
    const dropdownRef = useRef<HTMLDivElement>(null)

    const [user, setUser] = useState<User | null>(null)
    const [dropdownOpen, setDropdownOpen] = useState(false)

    useEffect(() => {
        fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/user`, {
            credentials: "include"
        })
        .then((res) => res.json())
        .then((data) => {
            if (data.User_Profile) {
                setUser(data.User_Profile)
            }
        })
        .catch((error) => console.error("Failed to fetch user:", error))
    }, [])

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (dropdownRef.current instanceof HTMLElement && !dropdownRef.current.contains(event.target as Node)) {
                setDropdownOpen(false)
            }
        }

        document.addEventListener("mousedown", handleClickOutside)
        return () => {
            document.removeEventListener("mousedown", handleClickOutside)
        }
    }, [])

    useEffect(() => {
        const handleLoginSuccess = (event: MessageEvent) => {
            if (event.origin !== process.env.NEXT_PUBLIC_BACKEND_API) {
                return
            }

            if (event.data === "login-success") {
                window.location.reload()
            }
        }

        window.addEventListener("message", handleLoginSuccess)
        return () => {
            window.removeEventListener("message", handleLoginSuccess)
        }
    }, [])

    const handleLoginClick = () => {
        const authWindow = window.open(
            `${process.env.NEXT_PUBLIC_BACKEND_API}/auth/google`,
            "Login",
            "width=960, height=540"
        )

        const checkLogin = setInterval(() => {
            if (!authWindow || authWindow.closed) {
                clearInterval(checkLogin)
                window.location.reload()
            }
        }, 1000)
    }

    const handleProfileClick = () => {
        setDropdownOpen(false)
        router.push("/user")
    }

    const handleLogout = async () => {
        await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/auth/logout`, {
            credentials: "include",
        })
        setUser(null)
        router.push("/")
    }

    const toggleDropdown = () => {
        setDropdownOpen(!dropdownOpen)
    }

    const handlePostClick = () => {
        if (user) {
            router.push("/f/create")
        } else {
            alert("You need to log in first to post your work!")
        }
    }

    return (
        <nav className="flex items-center justify-between p-4 bg-gray-900 text-white">
            <Link href="/" className="text-2xl font-bold hover:text-gray-300">
                Fictsu
            </Link>
            <div className="flex items-center space-x-4">
                {pathname !== "/f/create" && (
                    <button
                        onClick={handlePostClick}
                        className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
                    >
                        Post your work +
                    </button>
                )}
                {user ? (
                    <div className="relative" ref={dropdownRef}>
                        <button onClick={toggleDropdown} className="flex items-center space-x-2">
                            <img src={user.avatar_url} alt="Avatar" className="w-8 h-8 rounded-full" referrerPolicy="no-referrer" />
                            <span>{user.name}</span>
                        </button>
                        {dropdownOpen && (
                            <div className="absolute right-0 mt-2 w-48 bg-white text-black shadow-lg rounded-md overflow-hidden">
                                <button onClick={handleProfileClick} className="block w-full text-left px-4 py-2 hover:bg-gray-100">
                                    Profile
                                </button>
                                <button onClick={handleLogout} className="block w-full text-left px-4 py-2 hover:bg-gray-100">
                                    Logout
                                </button>
                            </div>
                        )}
                    </div>
                    ) : (
                    <button
                        onClick={handleLoginClick}
                        className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                    >
                        Login
                    </button>
                )}
            </div>
        </nav>
    )
}
