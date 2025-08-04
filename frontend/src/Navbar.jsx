export default function Navbar({ onLogin, username, showLogin }) {
  return (
    <nav className="w-full flex items-center justify-between px-6 py-4 absolute top-0 left-0 z-20">
      <div className="flex items-center">
        <span className="text-2xl font-bold text-white tracking-tight select-none">FinPal</span>
      </div>
      <div className="flex items-center">
        {username ? (
          <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-purple-500 to-pink-500 flex items-center justify-center text-white font-bold text-xl shadow-md select-none">
            {username[0].toUpperCase()}
          </div>
        ) : (
          showLogin && (
            <button
              className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-lg font-semibold shadow-md focus:ring-2 focus:ring-offset-2 focus:ring-pink-200 transition-all"
              onClick={onLogin}
            >
              Login
            </button>
          )
        )}
      </div>
    </nav>
  );
}
