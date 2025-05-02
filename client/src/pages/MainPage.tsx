import DocumentsList from '@/components/Documents/DocumentsList';
import Editor from '@/components/Editor/Editor';
import Header from '@/components/Header/Header';

function MainPage() {
    return (
        <div className="flex flex-col min-h-screen">
            <Header />
            <div className="flex flex-1">
                <DocumentsList />
                <Editor />
            </div>
        </div>
    );
}

export default MainPage;