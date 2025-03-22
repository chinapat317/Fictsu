import Link from "next/link"
import Image from "next/image"
import { Fiction } from "@/types/types"
import { notFound } from "next/navigation"
import FavoriteButton from "@/components/FavoriteButton"

async function getFiction(fiction_id: string): Promise<Fiction | null> {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}`, {
        cache: "no-store",
        credentials: "include",
    })

    if (!res.ok) {
        return null
    }

    const data = await res.json()
    return data.Fiction
}

export default async function FictionPage({ params }: { params: { fiction_id: string } }) {
    const { fiction_id } = await Promise.resolve(params)
    const fiction = await getFiction(fiction_id)
    if (!fiction) {
        notFound()
    }

    return (
        <div className="max-w-4xl mx-auto p-6">
            <div className="flex gap-6">
                <Image
                    src={fiction.cover || "/default-cover.png"}
                    alt={fiction.title}
                    width={200}
                    height={300}
                    className="rounded-lg"
                />
                <div>
                    <h1 className="text-3xl font-bold">{fiction.title}</h1>
                    <p className="text-gray-500">By {fiction.author}</p>
                    <FavoriteButton fiction_id={fiction.id} />
                    <p className="text-sm text-gray-400">Contributor: {fiction.contributor_name}</p>
                    <p className="mt-2 text-sm bg-gray-200 px-2 py-1 inline-block rounded-md">{fiction.status}</p>
                    <p className="mt-4">{fiction.synopsis}</p>
                </div>
            </div>
            {fiction.genres?.length > 0 && (
                <div className="mt-6">
                    <h2 className="text-xl font-semibold">Genres</h2>
                    <div className="flex flex-wrap gap-2 mt-2">
                        {fiction.genres.map((genre) => (
                            <span key={genre.id} className="bg-blue-200 px-3 py-1 rounded-lg text-sm">
                                {genre.genre_name}
                            </span>
                        ))}
                    </div>
                </div>
            )}
            {fiction.chapters?.length > 0 && (
                <div className="mt-6">
                    <h2 className="text-xl font-semibold">Chapters</h2>
                    <ul className="mt-2 list-disc list-inside">
                        {fiction.chapters.map((chapter) => (
                            <li key={chapter.id} className="flex items-center justify-between">
                                <Link href={`/f/${fiction.id}/${chapter.id}`} className="text-blue-500 hover:underline">
                                    {chapter.title}
                                </Link>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    )
}
