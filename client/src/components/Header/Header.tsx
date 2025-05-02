import { Link } from 'react-router-dom';
import { FiLogOut, FiFileText } from 'react-icons/fi';

const Header = () => {
  // Function to handle logout
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userId');
    window.location.href = '/login';
  };

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo/Brand */}
          <div className="flex-shrink-0 flex items-center">
            <Link to="/" className="text-xl font-serif text-gray-800 font-medium">
              PenFeel
            </Link>
          </div>
          
          {/* Navigation */}
          <nav className="flex items-center space-x-4">
            <Link 
              to="/" 
              className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium flex items-center"
            >
              <FiFileText className="mr-1" />
              Documents
            </Link>
            
            <div className="relative ml-3">
              <div className="flex">
                <button
                  type="button"
                  className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium flex items-center"
                  onClick={handleLogout}
                >
                  <FiLogOut className="mr-1" />
                  Log out
                </button>
              </div>
            </div>
          </nav>
        </div>
      </div>
    </header>
  );
};

export default Header; 