import { useState } from 'react';
import { FiPlus, FiFile } from 'react-icons/fi';
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

function DocumentsList() {
    const [error, setError] = useState<string | null>(null);
    const queryClient = useQueryClient();

    // Fetch documents query
    const { data, isLoading } = useQuery({
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
                
                // For other errors, try to get detailed error message from response
                try {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to fetch documents: ${response.status}`);
                } catch {
                    // If parsing JSON fails, use status text
                    throw new Error(`Failed to fetch documents: ${response.statusText || response.status}`);
                }
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
                // Try to get detailed error message from response
                try {
                    const errorData = await response.json();
                    throw new Error(errorData.error || `Failed to create document: ${response.status}`);
                } catch {
                    // If parsing JSON fails, use status text
                    throw new Error(`Failed to create document: ${response.statusText || response.status}`);
                }
            }
            
            return response.json();
        },
        onSuccess: () => {
            // Invalidate the documents query to trigger a refetch
            queryClient.invalidateQueries({ queryKey: ['documents'] });
        },
        onError: (error: Error) => {
            setError(error.message);
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
                    <div className="text-red-500 text-sm py-2">{error}</div>
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