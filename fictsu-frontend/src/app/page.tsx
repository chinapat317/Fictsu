import Link from "next/link"
import Image from "next/image"
import { Fiction } from "@/types/types"

async function getFictions(): Promise<Fiction[]> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}/f`, {
    cache: "no-store",
  })

  if (!res.ok) {
    throw new Error("Failed to fetch fictions")
  }

  return res.json()
}

export default async function HomePage() {
  const fictions = await getFictions()
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Fictsu Fictions</h1>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {fictions.map((fiction) => (
          <Link
            key={fiction.id}
            href={`/f/${fiction.id}`}
            className="block border rounded-lg overflow-hidden shadow-lg hover:shadow-xl transition"
          >
            <div className="relative w-full h-56">
              <Image
                src={fiction.cover || "/default-cover.png"}
                alt={fiction.title}
                layout="fill"
                objectFit="cover"
                className="rounded-t-lg"
              />
            </div>
            <div className="p-4">
              <h2 className="text-lg font-semibold">{fiction.title}</h2>
              <p className="text-sm text-gray-500">{fiction.subtitle}</p>
              <p className="text-sm text-gray-700">By {fiction.author}</p>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
