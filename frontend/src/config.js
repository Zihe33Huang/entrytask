// Define a new function called 'fetchWithJWT' that wraps the original 'fetch' function
export const fetchWithJWT = (input, init) => {
    // Retrieve the JWT from local storage
    const jwt = localStorage.getItem('jwt');
    // Create a new headers object or use the existing one in 'init' parameter
    const headers = init && init.headers ? new Headers(init.headers) : new Headers();

    // Set the JWT in the 'Authorization' header
    headers.set('Authorization', jwt);

    // Call the original 'fetch' function with the modified headers
    return fetch(input, { ...init, headers });
};

