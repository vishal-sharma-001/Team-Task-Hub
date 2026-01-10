import { useState, useCallback } from 'react';

export const useAsync = (asyncFunction, immediate = true) => {
  const [status, setStatus] = useState('idle');
  const [error, setError] = useState(null);
  const [data, setData] = useState(null);

  const execute = useCallback(async (...args) => {
    setStatus('pending');
    setData(null);
    setError(null);
    try {
      const response = await asyncFunction(...args);
      // Response is already unwrapped by the API interceptor
      setData(response);
      setStatus('success');
      return response;
    } catch (err) {
      setError(err.response?.data || err.message);
      setStatus('error');
      throw err;
    }
  }, [asyncFunction]);

  return { execute, status, error, data };
};

export const useForm = (initialValues, validations, onSubmit) => {
  const [values, setValues] = useState(initialValues);
  const [errors, setErrors] = useState({});
  const [touched, setTouched] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setValues((prev) => ({ ...prev, [name]: value }));
  };

  const handleBlur = (e) => {
    const { name } = e.target;
    setTouched((prev) => ({ ...prev, [name]: true }));

    // Validate on blur if validator exists
    if (validations && validations[name]) {
      const error = validations[name](values[name]);
      setErrors((prev) => ({ ...prev, [name]: error }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validate all fields
    const newErrors = {};
    if (validations) {
      Object.keys(validations).forEach((fieldName) => {
        const error = validations[fieldName](values[fieldName]);
        if (error) {
          newErrors[fieldName] = error;
        }
      });
    }

    setErrors(newErrors);
    setTouched(
      Object.keys(initialValues).reduce((acc, key) => {
        acc[key] = true;
        return acc;
      }, {})
    );

    // Only submit if no errors
    if (Object.keys(newErrors).length === 0) {
      setIsSubmitting(true);
      try {
        await onSubmit(values);
      } catch (err) {
        console.error('Form submission error:', err);
      } finally {
        setIsSubmitting(false);
      }
    }
  };

  return {
    values,
    errors,
    touched,
    isSubmitting,
    handleChange,
    handleBlur,
    handleSubmit,
    setValues,
  };
};

export const useLocalStorage = (key, initialValue) => {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const item = typeof window !== 'undefined' ? window.localStorage.getItem(key) : null;
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.error('Error reading from localStorage:', error);
      return initialValue;
    }
  });

  const setValue = (value) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      if (typeof window !== 'undefined') {
        window.localStorage.setItem(key, JSON.stringify(valueToStore));
      }
    } catch (error) {
      console.error('Error writing to localStorage:', error);
    }
  };

  return [storedValue, setValue];
};
