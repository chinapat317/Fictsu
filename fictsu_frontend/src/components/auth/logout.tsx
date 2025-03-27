import { useSession, signIn, signOut } from "next-auth/react";

export default function Logout() {
  const { data: session } = useSession();
  return (
    <div>
        <p>Welcome, {session.user?.name}</p>
        <button onClick={() => signOut()}>Sign Out</button>
    </div>
  );
};