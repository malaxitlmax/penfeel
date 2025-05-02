import { FiPlus, FiFile } from 'react-icons/fi';

function DocumentsList() {
    return (
        <div className="w-64 border-r border-gray-200 bg-gray-50 overflow-y-auto">
            <div className="p-4">
                <div className="flex items-center justify-between mb-4">
                    <h2 className="text-lg font-medium text-gray-700">Documents</h2>
                    <button 
                        className="p-1 rounded-full text-gray-600 hover:text-gray-900 hover:bg-gray-200"
                        title="Create new document"
                    >
                        <FiPlus className="h-5 w-5" />
                    </button>
                </div>
                
                <div className="space-y-2">
                    {/* Example documents - in a real app, these would come from an API */}
                    {[1, 2, 3].map((item) => (
                        <div 
                            key={item}
                            className="flex items-center p-2 rounded-md hover:bg-gray-200 cursor-pointer"
                        >
                            <FiFile className="mr-2 text-gray-600" />
                            <span className="text-gray-700">Document {item}</span>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

export default DocumentsList;