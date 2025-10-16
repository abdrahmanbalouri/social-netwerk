export async function middleware() {
    try {
        const response = await fetch("http://localhost:8080/api/me", {
            method: "GET",
            credentials: "include",
        });


        return response
    } catch (error) {

        throw error;
    }
}