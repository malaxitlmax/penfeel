import { useState, useEffect } from 'react';
import { FiPlus, FiFile, FiAlertCircle } from 'react-icons/fi';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { API_URL } from '@/App';

// Document interface to match our API
interface Document {
    id: string;
    title: string;
    content: string;
    user_id: string;
    created_at: string;
    updated_at: string;
}

// Error interface to capture detailed backend errors
interface DetailedError {
    message: string;
    details?: string;
    debug_info?: string;
}

function DocumentsList() {
    const [error, setError] = useState<DetailedError | null>(null);
    const queryClient = useQueryClient();

    // Helper to extract detailed error from API response
    const extractErrorDetails = async (response: Response, defaultMsg: string): Promise<DetailedError> => {
        try {
            const errorData = await response.json();
            return {
                message: errorData.error || defaultMsg,
                details: errorData.details,
                debug_info: errorData.debug_info
            };
        } catch {
            // If parsing JSON fails, use status text
            return {
                message: `${defaultMsg}: ${response.statusText || response.status}`
            };
        }
    };

    // Fetch documents query
    const { data, isLoading, error: queryError } = useQuery<Document[]>({
        queryKey: ['documents'],
        queryFn: async () => {
            const token = localStorage.getItem('token');
            
            const response = await fetch(`${API_URL}/documents`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (!response.ok) {
                // If response is 404, we want to handle it gracefully without retrying
                if (response.status === 404) {
                    return [];
                }
                
                // For other errors, get detailed error message
                const error = await extractErrorDetails(response, 'Failed to fetch documents');
                throw new Error(error.message);
            }
            
            const result = await response.json();
            return result.documents || [];
        },
        // Don't retry on 404 responses
        retry: (failureCount, error) => {
            // If the error message includes 404, don't retry
            if (error instanceof Error && error.message.includes('404')) {
                return false;
            }
            // Otherwise retry a few times (default behavior)
            return failureCount < 3;
        }
    });

    // Update error state when query errors occur
    useEffect(() => {
        if (queryError) {
            setError({ message: queryError.message });
        }
    }, [queryError]);

    // Create document mutation
    const createDocumentMutation = useMutation({
        mutationFn: async () => {
            const token = localStorage.getItem('token');
            
            const response = await fetch(`${API_URL}/documents`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ 
                    title: 'New Document', 
                    content: '' 
                })
            });
            
            if (!response.ok) {
                // Get detailed error info from response
                const errorDetails = await extractErrorDetails(response, 'Failed to create document');
                throw errorDetails;
            }
            
            return response.json();
        },
        onSuccess: () => {
            // Clear any previous errors
            setError(null);
            // Invalidate the documents query to trigger a refetch
            queryClient.invalidateQueries({ queryKey: ['documents'] });
        },
        onError: (err: unknown) => {
            // Handle the detailed error object
            if (typeof err === 'object' && err !== null) {
                setError(err as DetailedError);
            } else {
                setError({ message: String(err) });
            }
        }
    });

    const handleCreateDocument = () => {
        createDocumentMutation.mutate();
    };

    return (
        <div className="w-64 border-r border-gray-200 bg-gray-50 overflow-y-auto">
            <div className="p-4">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-lg font-medium text-gray-700">Documents</h2>
                    <button 
                        className="p-1 rounded-full text-gray-600 hover:text-gray-900 hover:bg-gray-200"
                        title="Create new document"
                        onClick={handleCreateDocument}
                    >
                        <FiPlus className="h-5 w-5" />
                    </button>
                </div>
                
                {isLoading ? (
                    <div className="flex justify-center py-4">
                        <span className="text-gray-500">Loading documents...</span>
                    </div>
                ) : error ? (
                    <div className="text-red-500 text-sm py-2 space-y-1">
                        <div className="flex items-start">
                            <FiAlertCircle className="mr-1 mt-0.5 flex-shrink-0" />
                            <span>{error.message}</span>
                        </div>
                        {error.details && (
                            <div className="text-xs text-red-400 ml-5">{error.details}</div>
                        )}
                    </div>
                ) : !data || data.length === 0 ? (
                    <div className="text-gray-500 text-sm py-2">
                        No documents found. Create your first document!
                    </div>
                ) : (
                    <div className="space-y-2">
                        {data.map((doc: Document) => (
                            <div 
                                key={doc.id}
                                className="flex items-center p-2 rounded-md hover:bg-gray-200 cursor-pointer"
                            >
                                <FiFile className="mr-2 text-gray-600" />
                                <span className="text-gray-700">{doc.title}</span>
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}

export default DocumentsList;