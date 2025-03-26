"use client"
import Link from "next/link"

export default function NavBar(){
    return(
    <nav className="flex items-center justify-between p-4 bg-gray-900 text-white">
        <Link href="/" className="text-2xl font-bold hover:text-gray-300">
                Fictsu
        </Link>
            <div className="flex items-center space-x-4">
                <button
                    onClick={() => {}}
                    className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700">
                    Post your work +
                </button>
                <button
                        onClick={() => {}}
                        className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
                        Login
                </button>
            </div>
    </nav>
    )
}