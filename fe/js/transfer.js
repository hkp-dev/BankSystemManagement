
document.addEventListener('DOMContentLoaded', function() {
    // Handle multi-step form
    const nextButtons = document.querySelectorAll('.next-step');
    const prevButtons = document.querySelectorAll('.prev-step');
    const steps = document.querySelectorAll('.step');
    
    // Next button click
    nextButtons.forEach(button => {
        button.addEventListener('click', function() {
            const currentStep = parseInt(this.getAttribute('data-step'));
            const prevStep = currentStep - 1;
            
            // Hide current step
            document.getElementById('step' + prevStep).style.display = 'none';
            
            // Show next step
            document.getElementById('step' + currentStep).style.display = 'block';
            
            // Update step indicators
            steps.forEach((step, index) => {
                if (index < currentStep) {
                    step.classList.add('completed');
                    step.classList.remove('active');
                }
                if (index === currentStep - 1) {
                    step.classList.add('active');
                }
            });
            
            // Update summary in step 3
            if (currentStep === 3) {
                updateSummary();
            }
            
            // Scroll to top of form
            document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
        });
    });
    
    // Previous button click
    prevButtons.forEach(button => {
        button.addEventListener('click', function() {
            const prevStep = parseInt(this.getAttribute('data-step'));
            const currentStep = prevStep + 1;
            
            // Hide current step
            document.getElementById('step' + currentStep).style.display = 'none';
            
            // Show previous step
            document.getElementById('step' + prevStep).style.display = 'block';
            
            // Update step indicators
            steps[currentStep - 1].classList.remove('active');
            steps[prevStep - 1].classList.add('active');
            
            // Scroll to top of form
            document.querySelector('.card').scrollIntoView({ behavior: 'smooth' });
        });
    });
    
    // Function to update summary in step 3
    function updateSummary() {
        // Get values from form
        const fromAccount = document.getElementById('fromAccount').options[document.getElementById('fromAccount').selectedIndex].text;
        const toAccount = document.getElementById('toAccount').value;
        const amount = document.getElementById('amount').value;
        const note = document.getElementById('note').value;
        
        // Update summary
        document.getElementById('summaryFromAccount').textContent = fromAccount;
        document.getElementById('summaryToAccount').textContent = toAccount;
        document.getElementById('summaryAmount').textContent = '$' + amount + '.00';
        document.getElementById('summaryNote').textContent = note || 'No description provided';
        
        // Update confirmation page
        document.getElementById('confirmAmount').textContent = '$' + amount + '.00';
        document.getElementById('confirmRecipient').textContent = 'Jane Doe';
    }
    
    // Format currency input
    document.getElementById('amount').addEventListener('input', function() {
        this.value = this.value.replace(/[^0-9]/g, '');
    });
    
    // Handle recipient selection
    const savedRecipients = [
        { name: 'Jane Doe', account: '****5678' },
        { name: 'John Smith', account: '****9012' },
        { name: 'Sarah Johnson', account: '****3456' }
    ];
    
    const recipientSelect = document.getElementById('recipientSelect');
    
    savedRecipients.forEach((recipient, index) => {
        const option = document.createElement('option');
        option.value = index;
        option.textContent = `${recipient.name} - ${recipient.account}`;
        recipientSelect.appendChild(option);
    });
    
    recipientSelect.addEventListener('change', function() {
        if (this.value !== '') {
            const selectedRecipient = savedRecipients[this.value];
            document.getElementById('toAccount').value = selectedRecipient.account;
            document.getElementById('recipientName').value = selectedRecipient.name;
        }
    });
});