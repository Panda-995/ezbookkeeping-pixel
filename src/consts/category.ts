import type { PresetCategory } from '@/core/category.ts';

// Presets follow the everyday order in which people think about spending:
// eat, wear, live, travel, play and learn. Tags intentionally stay empty so
// every household can build a lightweight taxonomy that matches its own life.
export const DEFAULT_EXPENSE_CATEGORIES: PresetCategory[] = [
    {
        name: 'Food & Drink',
        categoryIconId: '1',
        color: 'e15b36',
        subCategories: [
            { name: 'Food', categoryIconId: '2', color: 'e15b36' },
            { name: 'Drink', categoryIconId: '30', color: 'e15b36' },
            { name: 'Fruit & Snack', categoryIconId: '70', color: 'e15b36' }
        ]
    },
    {
        name: 'Clothing & Appearance',
        categoryIconId: '100',
        color: 'a85586',
        subCategories: [
            { name: 'Clothing', categoryIconId: '110', color: 'a85586' },
            { name: 'Cosmetic', categoryIconId: '180', color: 'a85586' }
        ]
    },
    {
        name: 'Housing & Houseware',
        categoryIconId: '200',
        color: '88683f',
        subCategories: [
            { name: 'Rent & Mortgage', categoryIconId: '290', color: '88683f' },
            { name: 'Utilities Expense', categoryIconId: '270', color: '88683f' },
            { name: 'Houseware', categoryIconId: '210', color: '88683f' },
            { name: 'Repairs & Maintenance', categoryIconId: '250', color: '88683f' }
        ]
    },
    {
        name: 'Transportation',
        categoryIconId: '300',
        color: '247d72',
        subCategories: [
            { name: 'Public Transit', categoryIconId: '310', color: '247d72' },
            { name: 'Taxi & Car Rental', categoryIconId: '320', color: '247d72' },
            { name: 'Personal Car Expense', categoryIconId: '330', color: '247d72' },
            { name: 'Train Tickets', categoryIconId: '370', color: '247d72' },
            { name: 'Airline Tickets', categoryIconId: '390', color: '247d72' }
        ]
    },
    {
        name: 'Entertainment',
        categoryIconId: '500',
        color: '6557a5',
        subCategories: [
            { name: 'Movies & Shows', categoryIconId: '550', color: '6557a5' },
            { name: 'Toys & Games', categoryIconId: '560', color: '6557a5' },
            { name: 'Sports & Fitness', categoryIconId: '510', color: '6557a5' },
            { name: 'Travelling', categoryIconId: '590', color: '6557a5' }
        ]
    },
    {
        name: 'Education & Studying',
        categoryIconId: '600',
        color: '8a721b',
        subCategories: [
            { name: 'Books & Newspaper & Magazines', categoryIconId: '610', color: '8a721b' },
            { name: 'Training Courses', categoryIconId: '660', color: '8a721b' },
            { name: 'Certification & Examination', categoryIconId: '680', color: '8a721b' }
        ]
    },
    {
        name: 'Communication',
        categoryIconId: '400',
        color: '2770a7',
        subCategories: [
            { name: 'Telephone Bill', categoryIconId: '420', color: '2770a7' },
            { name: 'Internet Bill', categoryIconId: '430', color: '2770a7' },
            { name: 'Subscriptions', categoryIconId: '570', color: '2770a7' }
        ]
    },
    {
        name: 'Medical & Healthcare',
        categoryIconId: '800',
        color: 'b34040',
        subCategories: [
            { name: 'Diagnosis & Treatment', categoryIconId: '840', color: 'b34040' },
            { name: 'Medications', categoryIconId: '860', color: 'b34040' }
        ]
    },
    {
        name: 'Gifts & Donations',
        categoryIconId: '700',
        color: '3a7d4f',
        subCategories: [
            { name: 'Gifts', categoryIconId: '710', color: '3a7d4f' },
            { name: 'Donations', categoryIconId: '780', color: '3a7d4f' }
        ]
    },
    {
        name: 'Finance & Insurance',
        categoryIconId: '900',
        color: 'b26b18',
        subCategories: [
            { name: 'Insurance Expense', categoryIconId: '950', color: 'b26b18' },
            { name: 'Tax Expense', categoryIconId: '910', color: 'b26b18' },
            { name: 'Service Charge', categoryIconId: '930', color: 'b26b18' }
        ]
    },
    {
        name: 'Miscellaneous',
        categoryIconId: '1000',
        color: '68736f',
        subCategories: [
            { name: 'Other Expense', categoryIconId: '1010', color: '68736f' }
        ]
    }
];

export const DEFAULT_INCOME_CATEGORIES: PresetCategory[] = [
    {
        name: 'Salary Income',
        categoryIconId: '2010',
        color: '247d5f',
        subCategories: [
            { name: 'Salary Income', categoryIconId: '2010', color: '247d5f' },
            { name: 'Bonus Income', categoryIconId: '2020', color: '247d5f' },
            { name: 'Overtime Pay', categoryIconId: '231', color: '247d5f' }
        ]
    },
    {
        name: 'Side Job Income',
        categoryIconId: '2080',
        color: '2770a7',
        subCategories: [
            { name: 'Side Job Income', categoryIconId: '2080', color: '2770a7' }
        ]
    },
    {
        name: 'Other Income',
        categoryIconId: '3010',
        color: '8a721b',
        subCategories: [
            { name: 'Investment Income', categoryIconId: '2100', color: '8a721b' },
            { name: 'Rental Income', categoryIconId: '290', color: '8a721b' },
            { name: 'Interest Income', categoryIconId: '970', color: '8a721b' },
            { name: 'Gift & Lucky Money', categoryIconId: '710', color: '8a721b' },
            { name: 'Winnings Income', categoryIconId: '564', color: '8a721b' },
            { name: 'Other Income', categoryIconId: '3010', color: '8a721b' }
        ]
    }
];

export const DEFAULT_TRANSFER_CATEGORIES: PresetCategory[] = [
    {
        name: 'General Transfer',
        categoryIconId: '4000',
        color: '247d72',
        subCategories: [
            { name: 'Bank Transfer', categoryIconId: '900', color: '247d72' },
            { name: 'Credit Card Repayment', categoryIconId: '980', color: '247d72' },
            { name: 'Deposits & Withdrawals', categoryIconId: '981', color: '247d72' }
        ]
    },
    {
        name: 'Loan & Debt',
        categoryIconId: '950',
        color: 'b26b18',
        subCategories: [
            { name: 'Borrowing Money', categoryIconId: '910', color: 'b26b18' },
            { name: 'Lending Money', categoryIconId: '290', color: 'b26b18' },
            { name: 'Repayment', categoryIconId: '930', color: 'b26b18' }
        ]
    },
    {
        name: 'Miscellaneous',
        categoryIconId: '1000',
        color: '68736f',
        subCategories: [
            { name: 'Reimbursement', categoryIconId: '920', color: '68736f' },
            { name: 'Other Transfer', categoryIconId: '4900', color: '68736f' }
        ]
    }
];
