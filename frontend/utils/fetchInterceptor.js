'use client'; // for Next.js 13+ app directory

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
export default function ClientFetchInterceptor() {
    const router = useRouter();
    useEffect(() => {
        const originalFetch = window.fetch.bind(window); // just copy

        window.fetch = async (...args) => {

            try {
                const response = await originalFetch(...args);
                if (response.status === 401) {
                    router.push('/login'); // rediret 
                    return response
                }
                return response;
            } catch (err) {
                console.error(err);
                throw err;
            }
        };
    }, []);

    return null;
}
