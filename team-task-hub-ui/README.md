# Team Task Hub - React Frontend

A modern, responsive React-based frontend for the Team Task Hub project management application.

## Features

- 🔐 **Authentication** - Secure login and signup with JWT tokens
- 📋 **Project Management** - Create, read, update, and delete projects
- ✅ **Task Management** - Full CRUD operations for tasks with status tracking
- 💬 **Comments** - Add comments and collaborate on tasks
- 📱 **Responsive Design** - Works seamlessly on desktop and mobile devices
- 🎨 **Modern UI** - Built with Tailwind CSS for beautiful styling
- ⚡ **Real-time Updates** - Instant task status changes and comments

## Tech Stack

- **Frontend Framework**: React 18.2.0
- **Routing**: React Router 6.18.0
- **HTTP Client**: Axios 1.6.0
- **Styling**: Tailwind CSS 3.3.0
- **Build Tool**: Vite 5.0.0
- **CSS Processing**: PostCSS with Tailwind

## Project Structure

```
src/
├── api/
│   └── client.js          # Axios instance with API endpoints
├── components/
│   ├── Navbar.jsx         # Navigation header
│   ├── ProtectedRoute.jsx # Route authentication guard
│   ├── Loading.jsx        # Loading spinner
│   ├── ErrorMessage.jsx   # Error display component
│   ├── ProjectForm.jsx    # Project create/edit form
│   ├── ProjectCard.jsx    # Project display card
│   ├── TaskForm.jsx       # Task create/edit form
│   ├── TaskCard.jsx       # Task display card
│   ├── CommentForm.jsx    # Comment input form
│   └── CommentList.jsx    # Comments display list
├── hooks/
│   └── useAsync.js        # Custom hooks (useAsync, useForm, useLocalStorage)
├── pages/
│   ├── Login.jsx          # Login page
│   ├── Signup.jsx         # Signup page
│   ├── Dashboard.jsx      # User dashboard
│   ├── Projects.jsx       # Projects list and management
│   └── TaskBoard.jsx      # Project tasks board
├── App.jsx                # Main app component with routing
├── main.jsx               # React DOM entry point
└── index.css              # Global styles and Tailwind config
```

## Installation

### Prerequisites

- Node.js 16.x or higher
- npm or yarn package manager
- Running backend API on `http://localhost:8080`

### Setup Steps

1. **Clone or navigate to the project directory**
   ```bash
   cd team-task-hub-ui
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Start the development server**
   ```bash
   npm run dev
   ```

   The application will be available at `http://localhost:3000`

4. **Build for production**
   ```bash
   npm run build
   ```

   Output files will be in the `dist/` directory.

## Available Scripts

- `npm run dev` - Start development server with hot reload
- `npm run build` - Build optimized production bundle
- `npm run preview` - Preview production build locally

## API Configuration

The frontend automatically proxies API requests to the backend:
- All `/api/*` requests are forwarded to `http://localhost:8080`
- Authentication token is automatically included in request headers
- 401 responses redirect to login

Configure the proxy in `vite.config.js` if your backend runs on a different port.

## Authentication

### Login Flow

1. User enters email and password on the Login page
2. Credentials are sent to `/api/auth/login`
3. Backend returns JWT token and user data
4. Token is stored in `localStorage` as `authToken`
5. User is redirected to Dashboard

### Protected Routes

All routes except `/login` and `/signup` require authentication. The `ProtectedRoute` component checks for a valid user in localStorage and redirects to login if not authenticated.

## Features Details

### Projects Page

- **View All Projects**: Lists all user's projects in a grid layout
- **Create Project**: Open modal form to create new project
- **Edit Project**: Modify project name and description
- **Delete Project**: Remove project with confirmation
- **View Tasks**: Navigate to task board for a specific project

### Task Board

- **View Tasks**: Display all tasks for a project
- **Create Task**: Add new task with title, description, priority
- **Edit Task**: Modify any task detail
- **Delete Task**: Remove task with confirmation
- **Change Status**: Update task status (OPEN → IN_PROGRESS → DONE)
- **Filter Tasks**: Filter by status and priority
- **Add Comments**: Comment on tasks for collaboration

### Task Card Features

- Priority indicator (Low/Medium/High) with color coding
- Status display and quick status change
- Assignee information
- Due date tracking
- Edit and delete actions

## Customization

### Tailwind CSS

Custom utilities and components are defined in `src/index.css`:

```css
.btn-primary       /* Primary action button */
.btn-secondary     /* Secondary action button */
.card              /* Card/container component */
.input             /* Form input styling */
.container         /* Max-width container */
```

Modify these in `index.css` to match your design system.

### Colors and Theme

Edit `tailwind.config.js` to customize:
- Color palette
- Spacing system
- Font families
- Border radius
- And more

## Form Validation

The `useForm` hook provides built-in validation support:

```javascript
const { values, errors, touched, handleChange, handleBlur, handleSubmit } = 
  useForm(
    { email: '', password: '' },
    {
      email: (value) => !value ? 'Email required' : '',
      password: (value) => !value ? 'Password required' : '',
    },
    async (formData) => {
      // Submit logic
    }
  );
```

## Error Handling

The application includes comprehensive error handling:

- **API Errors**: Displayed with `ErrorMessage` component
- **Form Validation**: Field-level error messages
- **Authentication**: Automatic redirect to login on 401
- **Network Errors**: User-friendly error display

## Development Tips

1. **Enable React DevTools** for browser debugging
2. **Use Vite's HMR** for fast hot module reloading during development
3. **Check Network tab** in DevTools to monitor API calls
4. **Use localStorage debugging** to inspect token and user data

## Deployment

### Build for Production

```bash
npm run build
```

This creates an optimized build in the `dist/` folder.

### Deployment Options

- **Vercel**: Connect GitHub repo for automatic deployments
- **Netlify**: Drag and drop `dist` folder or connect GitHub
- **GitHub Pages**: Build and serve static files
- **Docker**: Create a Docker image to serve with Node.js or nginx

### Environment Variables

For production deployments, you may need to configure the API base URL. Update `src/api/client.js`:

```javascript
const API_BASE = process.env.REACT_APP_API_URL || '/api';
```

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Troubleshooting

### API Connection Issues

- Ensure backend API is running on `http://localhost:8080`
- Check Network tab in DevTools for actual API URLs
- Verify CORS headers if hitting different domain

### Login Issues

- Clear `localStorage` and reload page
- Check browser console for error messages
- Verify credentials in backend database

### Build Issues

- Delete `node_modules` and `package-lock.json`, then `npm install`
- Clear Vite cache: `rm -rf .vite`
- Check Node.js version matches requirements

## Contributing

1. Create a feature branch
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## License

This project is part of Team Task Hub.

## Support

For issues, questions, or suggestions, please refer to the main project repository.
