// Main application logic
(function() {
    'use strict';

    // Application state
    let currentDocument = null;

    // DOM elements
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const browseBtn = document.getElementById('browseBtn');
    const fileInfo = document.getElementById('fileInfo');
    const uploadSection = document.getElementById('uploadSection');
    const generationSection = document.getElementById('generationSection');
    const resultsSection = document.getElementById('resultsSection');
    const generationForm = document.getElementById('generationForm');
    const loadingIndicator = document.getElementById('loadingIndicator');
    const errorMessage = document.getElementById('errorMessage');
    const generateAnotherBtn = document.getElementById('generateAnotherBtn');

    // Initialize app
    function init() {
        setupEventListeners();
        console.log('StudyForge initialized');
    }

    // Setup event listeners
    function setupEventListeners() {
        // Browse button
        browseBtn.addEventListener('click', () => fileInput.click());

        // File input change
        fileInput.addEventListener('change', handleFileSelect);

        // Drag and drop
        uploadArea.addEventListener('dragover', handleDragOver);
        uploadArea.addEventListener('dragleave', handleDragLeave);
        uploadArea.addEventListener('drop', handleDrop);

        // Generation form submit
        generationForm.addEventListener('submit', handleGenerationSubmit);

        // Generate another button
        generateAnotherBtn.addEventListener('click', showGenerationForm);
    }

    // Handle file selection
    async function handleFileSelect(event) {
        const file = event.target.files[0];
        if (file) {
            await uploadFile(file);
        }
    }

    // Handle drag over
    function handleDragOver(event) {
        event.preventDefault();
        uploadArea.classList.add('dragging');
    }

    // Handle drag leave
    function handleDragLeave(event) {
        event.preventDefault();
        uploadArea.classList.remove('dragging');
    }

    // Handle drop
    async function handleDrop(event) {
        event.preventDefault();
        uploadArea.classList.remove('dragging');

        const file = event.dataTransfer.files[0];
        if (file && file.type === 'application/pdf') {
            await uploadFile(file);
        } else {
            showError('Please drop a PDF file');
        }
    }

    // Upload file to server
    async function uploadFile(file) {
        // Validate file
        if (file.type !== 'application/pdf') {
            showError('Only PDF files are allowed');
            return;
        }

        const maxSize = 50 * 1024 * 1024; // 50MB
        if (file.size > maxSize) {
            showError('File size exceeds 50MB limit');
            return;
        }

        try {
            hideError();
            showLoading(uploadSection, 'Uploading PDF...');

            const result = await API.uploadPDF(file);
            currentDocument = result;

            // Update UI
            document.getElementById('fileName').textContent = result.filename;
            document.getElementById('pageCount').textContent = result.page_count;
            document.getElementById('fileSize').textContent = formatFileSize(result.file_size);

            fileInfo.style.display = 'block';
            hideLoading(uploadSection);

            // Show generation section
            generationSection.style.display = 'block';
            document.getElementById('pageEnd').max = result.page_count;
            document.getElementById('pageEnd').value = Math.min(10, result.page_count);

            console.log('Upload successful:', result);
        } catch (error) {
            hideLoading(uploadSection);
            showError('Upload failed: ' + error.message);
        }
    }

    // Handle generation form submit
    async function handleGenerationSubmit(event) {
        event.preventDefault();

        if (!currentDocument) {
            showError('No document uploaded');
            return;
        }

        const pageStart = parseInt(document.getElementById('pageStart').value);
        const pageEnd = parseInt(document.getElementById('pageEnd').value);
        const academicLevel = document.getElementById('academicLevel').value;

        // Validate page range
        if (pageStart < 1 || pageEnd < pageStart) {
            showError('Invalid page range');
            return;
        }

        if (pageEnd > currentDocument.page_count) {
            showError(`End page cannot exceed ${currentDocument.page_count}`);
            return;
        }

        try {
            hideError();
            loadingIndicator.style.display = 'block';
            generationForm.style.display = 'none';

            const result = await API.generateSummary(
                currentDocument.document_id,
                pageStart,
                pageEnd,
                academicLevel
            );

            // Display results
            displayResults(result, pageStart, pageEnd, academicLevel);

            console.log('Generation successful:', result);
        } catch (error) {
            showError('Generation failed: ' + error.message);
            loadingIndicator.style.display = 'none';
            generationForm.style.display = 'block';
        }
    }

    // Display results
    function displayResults(result, pageStart, pageEnd, academicLevel) {
        document.getElementById('resultPages').textContent = `${pageStart}-${pageEnd}`;
        document.getElementById('resultLevel').textContent = formatAcademicLevel(academicLevel);
        document.getElementById('resultTime').textContent = result.generation_time;
        document.getElementById('summaryContent').textContent = result.summary;

        loadingIndicator.style.display = 'none';
        generationSection.style.display = 'none';
        resultsSection.style.display = 'block';

        // Scroll to results
        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    // Show generation form again
    function showGenerationForm() {
        resultsSection.style.display = 'none';
        generationSection.style.display = 'block';
        generationForm.style.display = 'block';
        generationSection.scrollIntoView({ behavior: 'smooth' });
    }

    // Show loading indicator
    function showLoading(container, message) {
        const loading = document.createElement('div');
        loading.className = 'loading';
        loading.innerHTML = `
            <div class="spinner"></div>
            <p>${message}</p>
        `;
        loading.id = 'tempLoading';
        container.appendChild(loading);
    }

    // Hide loading indicator
    function hideLoading(container) {
        const loading = container.querySelector('#tempLoading');
        if (loading) {
            loading.remove();
        }
    }

    // Show error message
    function showError(message) {
        document.getElementById('errorText').textContent = message;
        errorMessage.style.display = 'block';

        // Auto-hide after 5 seconds
        setTimeout(() => {
            hideError();
        }, 5000);
    }

    // Hide error message
    function hideError() {
        errorMessage.style.display = 'none';
    }

    // Format file size
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    }

    // Format academic level
    function formatAcademicLevel(level) {
        const levels = {
            'high_school': 'High School',
            'undergraduate': 'Undergraduate',
            'graduate': 'Graduate'
        };
        return levels[level] || level;
    }

    // Start app when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
