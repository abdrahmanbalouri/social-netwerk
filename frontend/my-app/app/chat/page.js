"use client";
import { useRouter } from "next/navigation";

export default function ChatPage() {
    const router = useRouter();
    
    router.push("/chat/0");
}