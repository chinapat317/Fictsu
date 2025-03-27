"use client";

import { useSession, signIn} from "next-auth/react";
import Logout from "./logout";

export default function Login() {
  const { data: session } = useSession();

  const sendTokenToBackend = async () => {
    const res = await fetch("http://localhost:8080/api/f/c", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id_token: session?.id_token }),
    });

    const data = await res.json();
    console.log("Response from Go backend:", data);
  };

  return (
    <div>
      {!session ? (
        <button onClick={() => signIn("google")}>Sign in with Google</button>
      ) : (
        <Logout/>
      )}
    </div>
  );
}