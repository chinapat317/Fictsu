import Link from "next/link"
import { Chapter } from "@/types/types"
import { notFound } from "next/navigation"

async function getChapter(fiction_id: string, chapter_id: string): Promise<Chapter | null> {
    const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f/${fiction_id}/${chapter_id}`, {
        cache: "no-store",
    })
    if (!res.ok) {
        return null
    }

    const data = await res.json()
    return data
}

export default async function ChapterPage({ params }: { params: { fiction_id: string, chapter_id: string } }) {
    const { fiction_id, chapter_id } = await Promise.resolve(params)
    const chapter = await getChapter(fiction_id, chapter_id)
    if (!chapter) {
        notFound()
    }

    return (
        <div className="max-w-3xl mx-auto p-6">
            <h1 className="text-3xl font-bold">{chapter.title}</h1>
            <p className="text-gray-500 text-sm mt-2">
                Published on {new Date(chapter.created).toLocaleDateString()}
            </p>
            <div className="mt-6 text-lg leading-relaxed">
                {chapter.content.split("\n").map((paragraph, index) => (
                    <p key={index} className="mb-4">{paragraph}</p>
                ))}
            </div>
            <Link href={`/f/${chapter.fiction_id}`} className="text-blue-500 hover:underline mt-6 block">
                ← Back to Fiction
            </Link>
        </div>
    )
}
