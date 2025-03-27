// app/dashboard/page.tsx or pages/dashboard.tsx
"use client";

import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation"; // or 'next/router' if you're using pages/
import { useEffect } from "react";

export default function Authentiacate() {
  const { data: session, status } = useSession();
  const router = useRouter();
  useEffect(() => {
    if (status === "unauthenticated") {
      router.push("http://localhost:3000/home?login=false");
    }
  }, [status, router]);
  const loading = status === "loading";
  const authenticated = status === "authenticated";
  return { session, loading, authenticated };
}
