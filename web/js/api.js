// API module for communicating with the backend
const API = (function() {
    'use strict';

    const BASE_URL = '';

    // Upload PDF file
    async function uploadPDF(file) {
        const formData = new FormData();
        formData.append('file', file);

        try {
            const response = await fetch(`${BASE_URL}/api/documents/upload`, {
                method: 'POST',
                body: formData
            });

            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.error?.message || 'Upload failed');
            }

            return data.data;
        } catch (error) {
            console.error('Upload error:', error);
            throw error;
        }
    }

    // Generate study material
    async function generateSummary(documentId, pageStart, pageEnd, academicLevel) {
        try {
            const response = await fetch(`${BASE_URL}/api/study/generate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    document_id: documentId,
                    page_start: pageStart,
                    page_end: pageEnd,
                    material_type: 'summary',
                    academic_level: academicLevel
                })
            });

            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.error?.message || 'Generation failed');
            }

            return data.data;
        } catch (error) {
            console.error('Generation error:', error);
            throw error;
        }
    }

    // Get document information
    async function getDocument(documentId) {
        try {
            const response = await fetch(`${BASE_URL}/api/documents?id=${documentId}`);
            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.error?.message || 'Failed to get document');
            }

            return data.data;
        } catch (error) {
            console.error('Get document error:', error);
            throw error;
        }
    }

    // Health check
    async function healthCheck() {
        try {
            const response = await fetch(`${BASE_URL}/api/health`);
            const data = await response.json();
            return data.data;
        } catch (error) {
            console.error('Health check error:', error);
            return null;
        }
    }

    return {
        uploadPDF,
        generateSummary,
        getDocument,
        healthCheck
    };
})();
