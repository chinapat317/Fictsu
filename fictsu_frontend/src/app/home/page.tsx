"use client";

import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

export default function HomePage() {
  const searchParams = useSearchParams();
  const login = searchParams.get("login");

  useEffect(() => {
    if (login === "false") {
      alert("You must be logged in to view that page.");
    }
  }, [login]);

  return <div>This is the Home Page</div>;
}