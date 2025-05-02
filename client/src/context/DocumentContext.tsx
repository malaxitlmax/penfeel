import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import { API_URL } from '@/App';

// Document interface to match our API
export interface Document {
    id: string;
    title: string;
    content: string;
    user_id: string;
    created_at: string;
    updated_at: string;
}

interface DocumentContextType {
    selectedDocument: Document | null;
    documents: Document[];
    isLoading: boolean;
    error: Error | null;
    selectDocument: (documentId: string) => void;
    updateDocument: (id: string, content: string) => Promise<void>;
}

const DocumentContext = createContext<DocumentContextType | undefined>(undefined);

export function useDocumentContext() {
    const context = useContext(DocumentContext);
    if (context === undefined) {
        throw new Error('useDocumentContext must be used within a DocumentProvider');
    }
    return context;
}

interface DocumentProviderProps {
    children: ReactNode;
}

export function DocumentProvider({ children }: DocumentProviderProps) {
    const [selectedDocument, setSelectedDocument] = useState<Document | null>(null);
    const queryClient = useQueryClient();
    
    // Fetch documents query
    const { data: documents = [], isLoading, error } = useQuery<Document[]>({
        queryKey: ['documents'],
        queryFn: async () => {
            const token = localStorage.getItem('token');
            
            const response = await fetch(`${API_URL}/documents`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (!response.ok) {
                // If response is 404, we want to handle it gracefully
                if (response.status === 404) {
                    return [];
                }
                
                throw new Error(`Failed to fetch documents: ${response.statusText}`);
            }
            
            const result = await response.json();
            return result.documents || [];
        },
    });
    
    // Update document mutation
    const updateDocumentMutation = useMutation({
        mutationFn: async ({ id, content }: { id: string; content: string }) => {
            const token = localStorage.getItem('token');
            
            const response = await fetch(`${API_URL}/documents/${id}`, {
                method: 'PATCH',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ content })
            });
            
            if (!response.ok) {
                throw new Error(`Failed to update document: ${response.statusText}`);
            }
            
            return response.json();
        },
        onSuccess: () => {
            // Invalidate the documents query to trigger a refetch
            queryClient.invalidateQueries({ queryKey: ['documents'] });
        }
    });
    
    // Select first document by default when documents are loaded
    useEffect(() => {
        if (documents.length > 0 && !selectedDocument) {
            setSelectedDocument(documents[0]);
        }
    }, [documents, selectedDocument]);
    
    // Function to select a document
    const selectDocument = (documentId: string) => {
        const doc = documents.find(d => d.id === documentId);
        if (doc) {
            setSelectedDocument(doc);
        }
    };
    
    // Function to update document content
    const updateDocument = async (id: string, content: string) => {
        await updateDocumentMutation.mutateAsync({ id, content });
        
        // Update the selected document if it's the one that was changed
        if (selectedDocument?.id === id) {
            const updatedDoc = { ...selectedDocument, content };
            setSelectedDocument(updatedDoc);
        }
    };
    
    const value = {
        selectedDocument,
        documents,
        isLoading,
        error,
        selectDocument,
        updateDocument
    };
    
    return (
        <DocumentContext.Provider value={value}>
            {children}
        </DocumentContext.Provider>
    );
} 