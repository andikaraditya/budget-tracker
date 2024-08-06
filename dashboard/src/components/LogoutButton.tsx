"use client";

import { handleLogout } from "@/api/auth";
import { useRouter } from "next/navigation";

function LogoutButton() {
  const router = useRouter();
  function logout() {
    handleLogout();
  }
  return (
    <button
      onClick={logout}
      className="text-slate-700 font-semibold border-2 border-slate-700 rounded-lg px-4 py-2"
    >
      Logout
    </button>
  );
}

export default LogoutButton;
